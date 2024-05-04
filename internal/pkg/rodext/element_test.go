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
	element        *Element
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
}

func (ts *ElementSuite) TearDownTest() {
	ts.putPage()
	ts.cleanupBrowser()
}

func TestElementSuite(t *testing.T) {
	suite.Run(t, new(ElementSuite))
}

func (ts *ElementSuite) TestElement_Element() {

}
