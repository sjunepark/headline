package scraper

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"log/slog"
	"time"
)

// browser is a wrapper around rod.Browser
// It has a pagePool which you can get Pages from.
type browser struct {
	rodBrowser *rod.Browser
	pagePool   chan *page
}

// page is a wrapper around rod.Page
type page struct {
	rodPage *rod.Page
}

// newBrowser initializes a new browser. After initialization, it runs methods defined in browserOptions.
// A cleanup function is returned to close the browser and all pages in the pagePool.
func newBrowser(options browserOptions) (b *browser, cleanup func(), err error) {
	rodBrowser := rod.New()

	if options.debug {
		l, err := launcher.New().Headless(false).Devtools(true).Launch()
		if err != nil {
			return nil, nil, err
		}
		rodBrowser = rodBrowser.Trace(true).SlowMotion(2 * time.Second).ControlURL(l)
	}

	err = rodBrowser.Connect()
	if err != nil {
		return nil, nil, err
	}
	if options.noDefaultDevice {
		rodBrowser = rodBrowser.NoDefaultDevice()
	}
	if options.incognito {
		rodBrowser, err = rodBrowser.Incognito()
		if err != nil {
			return nil, nil, err
		}
	}

	numberOfPages := options.pagePoolSize
	pagePool := make(chan *page, numberOfPages)
	for i := 0; i < numberOfPages; i++ {
		pagePool <- nil
	}

	b = &browser{
		rodBrowser: rodBrowser,
		pagePool:   pagePool,
	}
	return b, b.cleanup, nil
}

// browserOptions holds options configurable while initializing a browser.
// They align with methods available to call on a rod.Browser.
type browserOptions struct {
	noDefaultDevice bool
	incognito       bool
	debug           bool
	pagePoolSize    int
}

// cleanup
// 1. Closes all Pages in pagePool
// 2. Closes the browser
func (b *browser) cleanup() {
	// Cannot use the range keyword here because it will deadlock
	for i := 0; i < cap(b.pagePool); i++ {
		p := <-b.pagePool
		if p == nil {
			continue
		}
		p.cleanup()
	}

	browserCloseErr := b.rodBrowser.Close()
	if browserCloseErr != nil {
		slog.Error("failed to close browser", "error", browserCloseErr)
	}
}

// page returns a new page from the pool.
// Always make sure to call putPage after using the page, or else a deadlock can occur.
// The returned page is thread-safe.
func (b *browser) page() (p *page, putPage func(), err error) {
	options := pageOptions{
		windowFullscreen: false,
	}

	putPageFactory := func(page *page) func() {
		return func() {
			fmt.Printf("Putting back page with address %p\n", page)
			b.pagePool <- page
		}
	}

	p = <-b.pagePool
	fmt.Printf("Got page with address %p\n", p)

	if p == nil {
		// We don't need to handle the cleanup since the page is to be reused.
		p, _, err = b.newPage(options)
		if err != nil {
			return nil, nil, err
		}
		fmt.Printf("Created new page with address %p\n", p)
		return p, putPageFactory(p), nil
	}

	fmt.Printf("Reusing page with address %p\n", p)
	return p, putPageFactory(p), nil
}

// newPage initializes a new page. After initialization, it runs methods defined in pageOptions.
func (b *browser) newPage(options pageOptions) (p *page, cleanup func(), err error) {
	opts := proto.TargetCreateTarget{}
	rodPage, err := b.rodBrowser.Page(opts)
	if err != nil {
		return nil, nil, err
	}

	if options.windowFullscreen {
		err = rodPage.SetWindow(&proto.BrowserBounds{
			WindowState: proto.BrowserWindowStateFullscreen,
		})
		if err != nil {
			return nil, nil, err
		}
	} else {
		err = rodPage.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
			Width:             1920,
			Height:            1080,
			DeviceScaleFactor: 1,
			Mobile:            false,
		})
		if err != nil {
			return nil, nil, err
		}
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
