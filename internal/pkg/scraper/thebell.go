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

func NewTheBellScraper() (*TheBellScraper, error) {
	options := browserOptions{
		noDefaultDevice: true,
		incognito:       true,
		debug:           false,
	}
	b, err := newBrowser(options)
	if err != nil {
		return nil, err
	}

	return &TheBellScraper{
		browser: b,
	}, nil
}

func (s *TheBellScraper) Cleanup() {
	s.browser.cleanup()
}

// Scrape should
// 1. Navigate to thebell
// 2. Get the list of articles to scrape for each keyword
// 3. For each article, get the article content
func (s *TheBellScraper) Scrape(keywords []string) ([]model.Article, error) {
	p, putPage, err := s.browser.page()
	defer putPage()
	if err != nil {
		return nil, err
	}
	err = p.navigate("https://thebell.co.kr/free/content/Article.asp?svccode=00/")
	if err != nil {
		return nil, err
	}
	return nil, nil
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
func cleanTheBellArticleUrl(u string) string {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return u
	}
	query := parsedUrl.Query()
	key := query.Get("key")
	if key == "" {
		return u
	}
	query = url.Values{"key": []string{key}}
	parsedUrl.RawQuery = query.Encode()
	u = parsedUrl.String()
	return u
}
