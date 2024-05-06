package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/constant"
	"github.com/sejunpark/headline/internal/pkg/rodext"
)

func init() {
	DefaultBrowserOptions = defaultBrowserOptions()
}

var DefaultBrowserOptions rodext.BrowserOptions

func defaultBrowserOptions() rodext.BrowserOptions {
	options := &rodext.BrowserOptions{
		Debug:           false,
		NoDefaultDevice: true,
		Incognito:       true,
		PagePoolSize:    constant.PAGE_POOL_SIZE,
	}
	return *options
}
