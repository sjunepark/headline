package rodext

import (
	"github.com/stretchr/testify/suite"
	"log/slog"
	"path/filepath"
)

type BaseBrowserSuite struct {
	suite.Suite
	Browser        *Browser
	CleanupBrowser func()
	PagePoolSize   int
	logLevel       slog.Level
}

func (ts *BaseBrowserSuite) SetupBaseBrowserTest() {
	ts.logLevel = slog.SetLogLoggerLevel(slog.LevelDebug)

	options := DefaultBrowserOptions
	b, cleanup, err := NewBrowser(options)
	ts.NoErrorf(err, "failed to initialize Browser: %v", err)

	ts.Browser = b
	ts.CleanupBrowser = cleanup
	ts.PagePoolSize = options.PagePoolSize
}

func (ts *BaseBrowserSuite) TearDownBaseBrowserTest() {
	ts.CleanupBrowser()
	slog.SetLogLoggerLevel(ts.logLevel)
}

type BasePageSuite struct {
	BaseBrowserSuite
	Page         *Page
	PutPage      func()
	WikipediaURL string
}

func (ts *BasePageSuite) SetupBasePageSuite() {
	relPath := "./testdata/WikiMockup.html"
	absPath, err := filepath.Abs(relPath)
	ts.NoErrorf(err, "failed to get absolute path: %v", err)
	ts.WikipediaURL = "file://" + absPath

	ts.SetupBaseBrowserTest()
}

func (ts *BasePageSuite) TearDownBasePageSuite() {
	ts.TearDownBaseBrowserTest()
}

func (ts *BasePageSuite) SetupBasePageSubTest() {
	p, putPage, err := ts.Browser.Page()
	ts.NoError(err, "failed to get Page")
	ts.Page = p
	ts.PutPage = putPage

	err = ts.Page.Navigate(ts.WikipediaURL)
	ts.NoError(err, "failed to navigate")
}

func (ts *BasePageSuite) TearDownBasePageSubTest() {
	ts.PutPage()
}
