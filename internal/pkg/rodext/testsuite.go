package rodext

import (
	"github.com/stretchr/testify/suite"
	"log/slog"
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
