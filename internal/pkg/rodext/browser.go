package rodext

import (
	"context"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"log/slog"
	"time"
)

// Browser is a wrapper around rod.Browser
// It has a pagePool which you can get Pages from.
type Browser struct {
	ctx        context.Context
	cancel     context.CancelFunc
	rodBrowser *rod.Browser
	pagePool   chan *page
}

// page is a wrapper around rod.Page
type page struct {
	rodPage *rod.Page
}

// NewBrowser initializes a new Browser. After initialization, it runs methods defined in BrowserOptions.
// A Cleanup function is returned to close the Browser and all pages in the pagePool.
func NewBrowser(options BrowserOptions) (b *Browser, cleanup func(), err error) {
	ctx, cancel := context.WithCancel(context.Background())

	rodBrowser := rod.New()

	if options.Debug {
		l, launchErr := launcher.New().Headless(false).Devtools(true).Launch()
		if launchErr != nil {
			cancel()
			return nil, nil, launchErr
		}
		rodBrowser = rodBrowser.Trace(true).SlowMotion(2 * time.Second).ControlURL(l)
	}

	err = rodBrowser.Connect()
	if err != nil {
		cancel()
		return nil, nil, err
	}
	if options.NoDefaultDevice {
		rodBrowser = rodBrowser.NoDefaultDevice()
	}
	if options.Incognito {
		rodBrowser, err = rodBrowser.Incognito()
		if err != nil {
			cancel()
			return nil, nil, err
		}
	}

	numberOfPages := options.PagePoolSize
	pagePool := make(chan *page, numberOfPages)
	// You have to fill the pagePool with nil values first.
	// By filling it with nil values, you can avoid using the select statement.
	// Also, if you don't fill it with nil values first, it's hard to properly clean them up.
	// There are lots of complexities to handle when a channel is not full.
	for i := 0; i < numberOfPages; i++ {
		pagePool <- nil
	}

	b = &Browser{
		ctx:        ctx,
		cancel:     cancel,
		rodBrowser: rodBrowser,
		pagePool:   pagePool,
	}
	return b, b.Cleanup, nil
}

// BrowserOptions holds options configurable while initializing a Browser.
// They align with methods available to call on a rod.Browser.
type BrowserOptions struct {
	NoDefaultDevice bool
	Incognito       bool
	Debug           bool
	PagePoolSize    int
}

// Cleanup
// 1. Closes all Pages in pagePool
// 2. Closes the Browser
func (b *Browser) Cleanup() {
	// Need to run cancel before closing pagPool to make sure no new pages are put back to the pool
	b.cancel()
	// Since you cannot designate specific producers for the pagePool, channel close should be handled by Cleanup.
	close(b.pagePool)
	for p := range b.pagePool {
		if p == nil {
			continue
		}
		p.cleanup()
	}

	browserCloseErr := b.rodBrowser.Close()
	if browserCloseErr != nil {
		slog.Error("failed to close Browser", "error", browserCloseErr)
	}
}

// page returns a new page from the pagePool.
// If the pagePool is empty, a new page is created.
// Always make sure to call putPage after using the page.
// The returned page is thread-safe.
func (b *Browser) page() (p *page, putPage func(), err error) {
	if b.ctx.Err() != nil {
		return nil, nil, b.ctx.Err()
	}

	p = <-b.pagePool

	if p != nil {
		slog.Debug("page(): got page from pool", "address", p)
		return p, putPageFactory(b.ctx, b.pagePool, p), nil
	}

	p, _, err = b.newPage(pageOptions{windowFullscreen: false})
	if err != nil {
		return nil, nil, err
	}
	slog.Debug("page(): no page in pool, created new page", "address", p)
	return p, putPageFactory(b.ctx, b.pagePool, p), nil
}

func putPageFactory(ctx context.Context, pagePool chan *page, page *page) func() {
	return func() {
		if ctx.Err() != nil {
			page.cleanup()
			slog.Debug("putPage: context cancelled, cleaning up page", "address", page)
			return
		}
		pagePool <- page
		slog.Debug("putPage: putting back page", "address", page)
	}
}

// newPage initializes a new page. After initialization, it runs methods defined in pageOptions.
func (b *Browser) newPage(options pageOptions) (p *page, cleanup func(), err error) {
	if b.ctx.Err() != nil {
		return nil, nil, b.ctx.Err()
	}

	opts := proto.TargetCreateTarget{}
	rodPage, err := b.rodBrowser.Page(opts)
	if err != nil {
		return nil, nil, err
	}

	err = setScreenSize(rodPage, options.windowFullscreen)
	if err != nil {
		return nil, nil, err
	}

	p = &page{rodPage: rodPage}
	return p, p.cleanup, nil
}

// pageOptions holds options configurable while initializing a page.
// They align with methods available to call on a rod.Page.
type pageOptions struct {
	windowFullscreen bool
}

// cleanup closes the page
func (p *page) cleanup() {
	err := p.rodPage.Close()
	if err != nil {
		slog.Error("failed to close page", "error", err)
	}
}

// navigate navigates the page to the given url.
// It waits for the NetworkAlmostIdle event before returning.
func (p *page) navigate(url string) error {
	var err error
	// A more conservative approach would be to wait for the onLoad event using the WaitLoad method.
	wait := p.rodPage.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)
	err = p.rodPage.Navigate(url)
	if err != nil {
		return err
	}
	wait()

	return nil
}

// element returns the rod.Element for the given selector.
func (p *page) element(selector string) (*rod.Element, error) {
	el, err := p.rodPage.Element(selector)
	if err != nil {
		return nil, err
	}
	return el, nil
}
