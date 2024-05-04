package rodext

import (
	"context"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"log/slog"
	"time"
)

// Browser is a wrapper around rod.Browser
// It has a pagePool which you can get Pages from.
type Browser struct {
	ctx        context.Context
	cancel     context.CancelFunc
	rodBrowser *rod.Browser
	pagePool   chan *Page
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
	pagePool := make(chan *Page, numberOfPages)
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

// Page returns a new Page from the pagePool.
// If the pagePool is empty, a new Page is created.
// Always make sure to call putPage after using the Page.
// The returned Page is thread-safe.
func (b *Browser) Page() (p *Page, putPage func(), err error) {
	if b.ctx.Err() != nil {
		return nil, nil, b.ctx.Err()
	}

	p = <-b.pagePool

	if p != nil {
		slog.Debug("Page(): got Page from pool", "address", p)
		return p, putPageFactory(b.ctx, b.pagePool, p), nil
	}

	p, _, err = newPage(b, pageOptions{windowFullscreen: false})
	if err != nil {
		return nil, nil, err
	}
	slog.Debug("Page(): no Page in pool, created new Page", "address", p)
	return p, putPageFactory(b.ctx, b.pagePool, p), nil
}

func putPageFactory(ctx context.Context, pagePool chan *Page, page *Page) func() {
	return func() {
		if ctx.Err() != nil {
			page.cleanup()
			slog.Debug("putPage: context cancelled, cleaning up Page", "address", page)
			return
		}
		pagePool <- page
		slog.Debug("putPage: putting back Page", "address", page)
	}
}
