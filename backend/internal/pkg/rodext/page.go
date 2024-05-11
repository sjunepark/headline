package rodext

import (
	"github.com/cockroachdb/errors"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"log/slog"
)

// Page is a wrapper around rod.Page
type Page struct {
	rodPage *rod.Page
}

// newPage initializes a new Page. After initialization, it runs methods defined in pageOptions.
func newPage(b *Browser, options pageOptions) (p *Page, cleanup func(), err error) {
	if b.ctx.Err() != nil {
		return nil, nil, errors.Wrap(b.ctx.Err(), "context cancelled in newPage")
	}

	opts := proto.TargetCreateTarget{}
	rodPage, err := b.rodBrowser.Page(opts)
	if err != nil {
		return nil, nil, errors.Wrap(err, "b.rodBrowser.Page() failed")
	}

	err = setScreenSize(rodPage, options.windowFullscreen)
	if err != nil {
		return nil, nil, errors.Wrap(err, "setScreenSize() failed")
	}

	p = &Page{rodPage: rodPage}
	return p, p.cleanup, nil
}

// newPageFactory returns a function(a closure) that returns a new Page.
// The user doesn't have to remember the reference to the Browser, as it is already captured in the closure.
func newPageFactory(b *Browser, options pageOptions) func() (*Page, func(), error) {
	return func() (*Page, func(), error) {
		return newPage(b, options)
	}
}

// pageOptions holds options configurable while initializing a Page.
// They align with methods available to call on a rod.Page.
type pageOptions struct {
	windowFullscreen bool
}

// cleanup closes the Page
func (p *Page) cleanup() {
	functionName := "Page.cleanup"
	err := p.rodPage.Close()
	slog.Debug("Page closed.", "function", functionName, "address", p)
	if err != nil {
		slog.Error("failed to close Page.", "function", functionName, "error", err)
	}
}

// Navigate navigates the Page to the given url.
// It returns a function that waits for the element to be loaded.
// Make sure to call this function to ensure the elements are loaded before applying selectors.
func (p *Page) Navigate(url string) (waitElement func(string) error, err error) {
	// This wait is enough in most cases.
	waitNavigation := p.rodPage.WaitNavigation(proto.PageLifecycleEventNameNetworkAlmostIdle)
	err = p.rodPage.Navigate(url)
	if err != nil {
		return nil, errors.Wrapf(err, "p.rodPage.Navigate(%s) failed", url)
	}
	waitNavigation()

	// Using WaitLoad is the safest way to ensure the page is fully loaded,
	// but can be very slow in many cases.
	waitElement = func(selector string) error {
		waitErr := p.rodPage.WaitElementsMoreThan(selector, 0)
		if waitErr != nil {
			return errors.Wrapf(waitErr, "p.rodPage.WaitElementsMoreThan(%s, 1) failed", selector)
		}
		return nil
	}

	return waitElement, nil
}

func (p *Page) Element(selector string) (*Element, error) {
	elements, err := p.rodPage.Elements(selector)
	if err != nil {
		html, _ := p.rodPage.HTML()
		return nil, errors.Wrapf(err, "p.rodPage.Elements(%s) failed, p.rodPage.HTML()=%s", selector, html)
	}

	if len(elements) > 1 {
		return nil, errors.Wrapf(MultipleElementsFoundError, "p.rodPage.Elements(%s) found multiple elements", selector)
	}

	element := elements.First()
	if element == nil {
		html, _ := p.rodPage.HTML()
		return nil, errors.Wrapf(ElementNotFoundError, "p.rodPage.Elements(%s) found no elements, p.rodPage.HTML()=%s", selector, html)
	}

	return &Element{rodElement: element}, nil
}

func (p *Page) Elements(selector string) ([]*Element, error) {
	elements, err := p.rodPage.Elements(selector)
	if err != nil {
		html, _ := p.rodPage.HTML()
		return nil, errors.Wrapf(err, "p.rodPage.Elements(%s) failed for html %s", selector, html)
	}
	return newElements(elements), nil
}

func (p *Page) HTML() (string, error) {
	return p.rodPage.HTML()
}
