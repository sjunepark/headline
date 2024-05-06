package rodext

import (
	"errors"
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
		return nil, nil, errors.New("cancelled context")
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
	err := p.rodPage.Close()
	if err != nil {
		slog.Error("failed to close Page", "error", err)
	}
}

// Navigate navigates the Page to the given url.
// It waits for the NetworkAlmostIdle event before returning.
func (p *Page) Navigate(url string) error {
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

var MultipleElementsFoundError = errors.New("multiple elements found")
var NotFoundError = errors.New("element not found")

func (p *Page) Element(selector string) (*Element, error) {
	elements, err := p.rodPage.Elements(selector)
	if err != nil {
		return nil, err
	}

	if len(elements) > 1 {
		return nil, MultipleElementsFoundError
	}

	element := elements.First()
	if element == nil {
		return nil, NotFoundError
	}

	return &Element{rodElement: element}, nil
}

func (p *Page) Elements(selector string) ([]*Element, error) {
	elements, err := p.rodPage.Elements(selector)
	if err != nil {
		return nil, err
	}
	return newElements(elements), nil
}
