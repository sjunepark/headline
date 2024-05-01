package scraper

import (
	"fmt"
	"github.com/sejunpark/headline/internal/pkg/model"
	"net/url"
)

type TheBellScraper struct {
	sourceScraper
	browser *browser
}

func NewTheBellScraper() (theBellScraper *TheBellScraper, cleanup func(), err error) {
	options := browserOptions{
		noDefaultDevice: true,
		incognito:       true,
		debug:           false,
	}
	b, browserCleanup, err := newBrowser(options)
	if err != nil {
		return nil, nil, err
	}

	theBellScraper = &TheBellScraper{
		browser: b,
	}
	return theBellScraper, browserCleanup, nil
}

func (s *TheBellScraper) Cleanup() {
	s.browser.cleanup()
}

func (s *TheBellScraper) getArticleUrls(keyword string) (<-chan []url.URL, error) {
	baseUrl := fmt.Sprintf("https://thebell.co.kr/free/content/Search.asp?page=1&period=360&part=A&keyword=%s", keyword)
	fmt.Printf("baseUrl: %s\n", baseUrl)
	// todo: implement
	return nil, nil
}

func (s *TheBellScraper) getArticle(url string) (model.Article, error) {
	return model.Article{}, nil
}

// cleanTheBellArticleUrl removes unnecessary query parameters from thebell article url,
// leaving only the 'key' parameter
func cleanTheBellArticleUrl(u string) (string, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	query := parsedUrl.Query()
	key := query.Get("key")
	if key == "" {
		return "", fmt.Errorf("parameter 'key' not found in url: %s", u)
	}
	query = url.Values{"key": []string{key}}
	parsedUrl.RawQuery = query.Encode()
	u = parsedUrl.String()
	return u, nil
}
