package scraper

type TheBellScraper struct {
	browser *browser
}

func NewTheBellScraper() (*TheBellScraper, error) {
	browserOptions := browserOptions{
		noDefaultDevice: true,
		incognito:       true,
		debug:           false,
	}
	browser, err := newBrowser(browserOptions)
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
