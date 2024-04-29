package scraper

import "github.com/sejunpark/headline/internal/pkg/rodext"

type TheBellScraper struct {
	browser *rodext.Browser
}

func NewTheBellScraper() (*TheBellScraper, error) {
	browserOptions := rodext.BrowserOptions{
		NoDefaultDevice: true,
		Incognito:       true,
		Debug:           false,
	}
	browser, err := rodext.NewBrowser(browserOptions)
	if err != nil {
		return nil, err
	}

	return &TheBellScraper{
		browser: browser,
	}, nil
}

func (s *TheBellScraper) Scrape() error {
	// todo: implement scraping logic
	return nil
}
