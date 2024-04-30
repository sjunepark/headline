package scraper

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// Page is a wrapper around rod.Page
type Page struct {
	rodPage *rod.Page
}

// newPage initializes a new Page. After initialization, it runs methods defined in pageOptions.
func newPage(b *Browser, options pageOptions) (*Page, error) {
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

func newPageFactory(b *Browser, options pageOptions) func() (*Page, error) {
	return func() (*Page, error) {
		return newPage(b, options)
	}
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
// It waits until the Page is stable for 1 second, and sets the window state to fullscreen.
func (p *Page) navigate(url string) error {
	var err error
	wait := p.rodPage.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)
	err = p.rodPage.Navigate(url)
	if err != nil {
		return err
	}
	wait()

	return nil
}

func (p *Page) Element(selector string) (*rod.Element, error) {
	el, err := p.rodPage.Element(selector)
	if err != nil {
		return nil, err
	}
	return el, nil
}

// pagePool is a channel of *Page
type pagePool chan *Page

func newPagePool(limit int) pagePool {
	pp := make(chan *Page, limit)
	for i := 0; i < limit; i++ {
		pp <- nil
	}
	return pp
}

// Get returns a Page from the pool.
// If the pool is empty, it creates a new Page according to the create function.
func (pp pagePool) get(create func() (*Page, error)) (*Page, error) {
	var err error
	p := <-pp
	if p == nil {
		p, err = create()
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

// put returns the used Page back to the pool.
func (pp pagePool) put(p *Page) {
	pp <- p
}

// cleanup iterates over all pages in the pool and runs the cleanup method for each Page.
func (pp pagePool) cleanup() error {
	for i := 0; i < cap(pp); i++ {
		p := <-pp
		if p != nil {
			err := p.cleanup()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
