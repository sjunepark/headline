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
	pagePool   *pagePool
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
	pp := newPagePool(ctx, numberOfPages)

	b = &Browser{
		ctx:        ctx,
		cancel:     cancel,
		rodBrowser: rodBrowser,
		pagePool:   pp,
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

	b.pagePool.cleanup()

	browserCloseErr := b.rodBrowser.Close()
	if browserCloseErr != nil {
		slog.Error("failed to close Browser", "error", browserCloseErr)
	}
}

// Page returns a new Page from the pool.
// Make sure to call putPage after using the Page, to put it back for reuse.
func (b *Browser) Page() (p *Page, putPage func(), err error) {
	options := pageOptions{
		windowFullscreen: false,
	}
	newPageFunc := newPageFactory(b, options)
	return b.pagePool.Get(newPageFunc)
}
