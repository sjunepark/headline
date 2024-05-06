package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/constant"
	"github.com/sejunpark/headline/internal/pkg/rodext"
)

var browserOptions = rodext.BrowserOptions{
	Debug:           false,
	NoDefaultDevice: true,
	Incognito:       true,
	PagePoolSize:    constant.PAGE_POOL_SIZE,
}
