package scraper

type TheBellScraper struct {
	browser *Browser
}

func NewTheBellScraper() (*TheBellScraper, error) {
	browserOptions := BrowserOptions{
		NoDefaultDevice: true,
		Incognito:       true,
		Debug:           false,
	}
	browser, err := NewBrowser(browserOptions)
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
