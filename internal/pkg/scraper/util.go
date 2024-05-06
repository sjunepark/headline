package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/constant"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"os"
)

func init() {
	DefaultBrowserOptions = defaultBrowserOptions()
}

var DefaultBrowserOptions rodext.BrowserOptions

func defaultBrowserOptions() rodext.BrowserOptions {
	options := DefaultBrowserOptions
	if os.Getenv("CI") == "true" {
		options = rodext.BrowserOptions{
			Debug:           true,
			NoDefaultDevice: true,
			Incognito:       true,
			PagePoolSize:    constant.PAGE_POOL_SIZE,
		}
	} else {
		options = rodext.BrowserOptions{
			Debug:           false,
			NoDefaultDevice: true,
			Incognito:       true,
			PagePoolSize:    constant.PAGE_POOL_SIZE,
		}
	}
	return options
}
