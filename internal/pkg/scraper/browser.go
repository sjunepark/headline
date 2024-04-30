package scraper

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"time"
)

// Browser is a wrapper around rod.Browser
// It has a pagePool which you can get Pages from.
type Browser struct {
	rodBrowser *rod.Browser
	pagePool   chan *Page
}

// Page is a wrapper around rod.Page
type Page struct {
	rodPage *rod.Page
}

// NewBrowser initializes a new Browser. After initialization, it runs methods defined in BrowserOptions.
func NewBrowser(options BrowserOptions) (*Browser, error) {
	var err error
	b := rod.New()

	if options.Debug {
		l, err := launcher.New().Headless(false).Devtools(true).Launch()
		if err != nil {
			return nil, err
		}
		b = b.Trace(true).SlowMotion(2 * time.Second).ControlURL(l)
	}

	err = b.Connect()
	if err != nil {
		return nil, err
	}
	if options.NoDefaultDevice {
		b = b.NoDefaultDevice()
	}
	if options.Incognito {
		b, err = b.Incognito()
		if err != nil {
			return nil, err
		}
	}

	numberOfPages := options.PagePoolSize
	pagePool := make(chan *Page, numberOfPages)
	for i := 0; i < numberOfPages; i++ {
		pagePool <- nil
	}

	return &Browser{
		rodBrowser: b,
		pagePool:   pagePool,
	}, nil
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
func (b *Browser) Cleanup() error {
	// Cannot use the range keyword here because it will deadlock
	for i := 0; i < cap(b.pagePool); i++ {
		page := <-b.pagePool
		if page == nil {
			continue
		}
		err := page.cleanup()
		if err != nil {
			return err
		}
	}

	err := b.rodBrowser.Close()
	if err != nil {
		return err
	}

	return nil
}

// Page returns a new Page from the pool.
// Always make sure to call putPage after using the Page, or else a deadlock can occur.
// The returned Page is thread-safe.
func (b *Browser) Page() (page *Page, putPage func(), err error) {
	options := pageOptions{
		WindowFullscreen: false,
	}

	putPageFactory := func(page *Page) func() {
		return func() {
			fmt.Printf("Putting back page with address %p\n", page)
			b.pagePool <- page
		}
	}

	page = <-b.pagePool
	fmt.Printf("Got page with address %p\n", page)

	if page == nil {
		page, err = b.newPage(options)
		if err != nil {
			return nil, nil, err
		}
		fmt.Printf("Created new page with address %p\n", page)
		return page, putPageFactory(page), nil
	}

	fmt.Printf("Reusing page with address %p\n", page)
	return page, putPageFactory(page), nil
}

// newPage initializes a new Page. After initialization, it runs methods defined in pageOptions.
func (b *Browser) newPage(options pageOptions) (*Page, error) {
	opts := proto.TargetCreateTarget{}
	p, err := b.rodBrowser.Page(opts)
	if err != nil {
		return nil, err
	}

	if options.WindowFullscreen {
		err = p.SetWindow(&proto.BrowserBounds{
			WindowState: proto.BrowserWindowStateFullscreen,
		})
		if err != nil {
			return nil, err
		}
	} else {
		err = p.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
			Width:             1920,
			Height:            1080,
			DeviceScaleFactor: 1,
			Mobile:            false,
		})
		if err != nil {
			return nil, err
		}
	}

	return &Page{rodPage: p}, nil
}

// pageOptions holds options configurable while initializing a Page.
// They align with methods available to call on a rod.Page.
type pageOptions struct {
	WindowFullscreen bool
}

// cleanup closes the Page
func (p *Page) cleanup() error {
	err := p.rodPage.Close()
	if err != nil {
		return err
	}
	return nil
}

// navigate navigates the Page to the given url.
// It waits for the NetworkAlmostIdle event before returning.
func (p *Page) navigate(url string) error {
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

// Element returns the rod.Element for the given selector.
func (p *Page) Element(selector string) (*rod.Element, error) {
	el, err := p.rodPage.Element(selector)
	if err != nil {
		return nil, err
	}
	return el, nil
}
