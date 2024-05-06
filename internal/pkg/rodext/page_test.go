package rodext

import (
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type PageSuite struct {
	suite.Suite
	browser        *Browser
	page           *Page
	putPage        func()
	cleanupBrowser func()
	wikipediaURL   string
}

func (ts *PageSuite) SetupSuite() {
	relPath := "./testdata/WikiMockup.html"
	absPath, err := filepath.Abs(relPath)
	ts.NoErrorf(err, "failed to get absolute path: %v", err)
	ts.wikipediaURL = "file://" + absPath
}

func (ts *PageSuite) SetupTest() {
	options := BrowserOptions{
		NoDefaultDevice: true,
		Incognito:       true,
		Debug:           false,
		PagePoolSize:    16,
	}
	b, cleanupBrowser, err := NewBrowser(options)
	ts.NoErrorf(err, "failed to initialize Browser: %v", err)
	ts.browser = b
	ts.cleanupBrowser = cleanupBrowser
}

func (ts *PageSuite) TearDownTest() {
	ts.cleanupBrowser()
}

func (ts *PageSuite) SetupSubTest() {
	p, putPage, err := ts.browser.Page()
	ts.NoErrorf(err, "failed to get Page: %v", err)
	ts.page = p
	ts.putPage = putPage

	err = ts.page.Navigate(ts.wikipediaURL)
	ts.NoErrorf(err, "failed to navigate: %v", err)
}

func (ts *PageSuite) TearDownSubTest() {
	ts.putPage()
}

func TestPageSuite(t *testing.T) {
	suite.Run(t, new(PageSuite))
}

func (ts *PageSuite) TestPage_cleanup() {
	ts.Run("Page should be cleaned up after cleanup is called", func() {
		page := ts.page
		_, err := page.rodPage.Info()
		ts.NoErrorf(err, "failed to get Page info: %v", err)

		page.cleanup()

		_, err = page.rodPage.Info()
		ts.Error(err, "shouldn't be able to get Page info after Cleanup")
	})
}

func (ts *PageSuite) TestPage_Element() {
	ts.Run("Element should return the element if it exists", func() {
		el, err := ts.page.Element("h1")
		ts.NoErrorf(err, "failed to get element: %v", err)
		ts.NotNilf(el, "element should not be nil: %v", el)
	})

	ts.Run("Element should return an error when multiple elements are found", func() {
		el, err := ts.page.Element("a")
		ts.ErrorIsf(err, MultipleElementsFoundError, "should return MultipleElementsFoundError when multiple elements are found, got: %v", err)
		ts.Nilf(el, "element should be nil: %v", el)
	})

	ts.Run("Element should return an error when no elements are found", func() {
		el, err := ts.page.Element(".nonexistent-element")
		ts.ErrorIsf(err, ElementNotFoundError, "should return ElementNotFoundError when no elements are found, got: %v", err)
		ts.Nil(el, "element should be nil: %v", el)
	})
}

func (ts *PageSuite) TestPage_Elements() {
	ts.Run("Elements should return the elements if they exist", func() {
		els, err := ts.page.Elements("a")
		ts.NoErrorf(err, "failed to get elements: %v", err)
		ts.NotEmptyf(els, "elements should not be empty: %v", els)
	})

	ts.Run("Elements should return an empty slice when no elements are found", func() {
		els, err := ts.page.Elements(".nonexistent-element")
		ts.NoErrorf(err, "failed to get elements: %v", err)
		ts.Empty(els, "elements should be empty: %v", els)
	})
}
