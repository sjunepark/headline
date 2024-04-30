package scraper

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"time"
)

// Browser is a wrapper around rod.Browser
type Browser struct {
	rodBrowser *rod.Browser
	pagePool   pagePool
}

// BrowserOptions holds options configurable while initializing a Browser.
// They align with methods available to call on a rod.Browser.
type BrowserOptions struct {
	NoDefaultDevice bool
	Incognito       bool
	Debug           bool
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

	// todo: need to determine how many pages to pool
	numberOfPages := 16
	pagePool := newPagePool(numberOfPages)

	return &Browser{
		rodBrowser: b,
		pagePool:   pagePool,
	}, nil
}

func (b *Browser) Cleanup() error {
	err := b.pagePool.cleanup()
	if err != nil {
		return err
	}
	return nil
}

// Page returns a new Page from the pool.
// Always make sure to call PutPage after using the Page, or else a deadlock will occur.
func (b *Browser) Page() (*Page, error) {
	options := pageOptions{
		WindowFullscreen: false,
	}
	create := newPageFactory(b, options)

	p, err := b.pagePool.get(create)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (b *Browser) PutPage(page *Page) {
	b.pagePool.put(page)
}
