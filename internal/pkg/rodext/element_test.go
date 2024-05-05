package rodext

import (
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type ElementSuite struct {
	suite.Suite
	browser        *Browser
	page           *Page
	ulElement      *Element
	putPage        func()
	cleanupBrowser func()
	wikipediaURL   string
}

func (ts *ElementSuite) SetupSuite() {
	relPath := "./testdata/WikiMockup.html"
	absPath, err := filepath.Abs(relPath)
	ts.NoErrorf(err, "failed to get absolute path: %v", err)
	ts.wikipediaURL = "file://" + absPath
}

func (ts *ElementSuite) SetupTest() {
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

	p, putPage, err := ts.browser.Page()
	ts.NoErrorf(err, "failed to get Page: %v", err)
	ts.page = p
	ts.putPage = putPage

	err = ts.page.navigate(ts.wikipediaURL)
	ts.NoErrorf(err, "failed to navigate: %v", err)

	ts.ulElement, err = ts.page.Element("ul")
	ts.NoErrorf(err, "failed to get Element: %v", err)
}

func (ts *ElementSuite) TearDownTest() {
	ts.putPage()
	ts.cleanupBrowser()
}

func TestElementSuite(t *testing.T) {
	suite.Run(t, new(ElementSuite))
}

func (ts *ElementSuite) TestElement_Element() {
	ts.Run("should return an Element with the correct selector", func() {
		aElement, err := ts.ulElement.Element("a[href='#section1']")
		ts.NoErrorf(err, "failed to get Element: %v", err)
		ts.Truef(isElement(aElement), "returned Element is not an Element")
		ts.NotNilf(aElement, "returned Element is nil")
	})

	ts.Run("should error if the selector is invalid", func() {
		el, err := ts.ulElement.Element("invalid")
		ts.ErrorIsf(err, ElementNotFoundError, "expected ElementNotFoundError but got %v", err)
		ts.Nilf(el, "expected Element to be nil but got %v", el)
	})

	ts.Run("should error when multiple elements are found", func() {
		el, err := ts.ulElement.Element("li")
		ts.ErrorIsf(err, MultipleElementsFoundError, "expected MultipleElementsFoundError but got %v", err)
		ts.Nilf(el, "expected Element to be nil but got %v", el)
	})
}

func (ts *ElementSuite) TestElement_Elements() {
	ts.Run("should return a slice of Elements with the correct selector", func() {
		liElements, err := ts.ulElement.Elements("li")
		ts.NoErrorf(err, "failed to get Elements: %v", err)
		ts.Truef(len(liElements) > 0, "expected slice of Elements to have length > 0 but got %v", len(liElements))
		for _, el := range liElements {
			ts.Truef(isElement(el), "returned Element is not an Element")
		}
	})

	ts.Run("should error if the selector is invalid", func() {
		els, err := ts.ulElement.Elements("invalid")
		ts.ErrorIsf(err, ElementNotFoundError, "expected ElementNotFoundError but got %v", err)
		ts.Nilf(els, "expected Elements to be nil but got %v", els)
	})
}

func isElement(toCheck any) bool {
	switch toCheck.(type) {
	case *Element:
		return true
	}
	return false
}
