package scraper

import (
	"fmt"
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"net/url"
	"time"
)

type ThebellScraper struct {
	browser *rodext.Browser
}

func NewThebellScraper() (s *ThebellScraper, cleanup func(), err error) {
	options := rodext.BrowserOptions{
		NoDefaultDevice: true,
		Incognito:       true,
		Debug:           false,
		PagePoolSize:    16,
	}
	b, browserCleanup, err := rodext.NewBrowser(options)
	if err != nil {
		return nil, nil, err
	}

	s = &ThebellScraper{
		browser: b,
	}
	return s, browserCleanup, nil
}

func (s *ThebellScraper) cleanup() {
	s.browser.Cleanup()
}

type articleMetadata struct {
	title           string
	createdDateTime time.Time // format: "2006-01-02 15:04:05"
	relativeUrl     *url.URL
	summary         string
}

func (s *ThebellScraper) fetchArticleList(keyword string) (<-chan *model.Article, error) {
	baseUrl := fmt.Sprintf("https://thebell.co.kr/free/content/Search.asp?page=1&period=360&part=A&keyword=%s", keyword)
	fmt.Printf("baseUrl: %s\n", baseUrl)

	page, putPage, err := s.browser.Page()
	if err != nil {
		return nil, err
	}
	defer putPage()

	articleElementsToScrape, err := page.Elements("div.listBox>ul>li>dl")
	if err != nil {
		return nil, err
	}
	for _, articleElementToScrape := range articleElementsToScrape {
		// todo: implement
		articleListItem, parseErr := s.parseArticleListItem(articleElementToScrape)
		fmt.Printf("articleListItem: %+v\n", articleListItem)
		if parseErr != nil {
			continue
		}
	}

	return nil, nil
}

func (s *ThebellScraper) fetchArticle(url *url.URL) (*model.Article, error) {
	// todo: implement
	return &model.Article{}, nil
}

func (s *ThebellScraper) parseArticleListItem(el *rodext.Element) (*model.Article, error) {
	a, elParseErr := el.Element("a")
	if elParseErr != nil {
		return nil, elParseErr
	}

	href := a.Attribute("href")
	relativeUrl, urlParseErr := url.Parse(href)
	if urlParseErr != nil {
		return nil, urlParseErr
	}
	title := a.Attribute("title")

	datetimeEl, datetimeElParseErr := el.Element("dd>.date")
	if datetimeElParseErr != nil {
		return nil, datetimeElParseErr
	}

	datetime, err := parseThebellKoreanDatetime(datetimeEl.Text())
	if err != nil {
		return nil, err
	}

	return &model.Article{
		Title:           title,
		CreatedDateTime: datetime,
		UpdatedDateTime: datetime,
		Source:          "",
		Url:             relativeUrl,
		Summary:         "",
		Content:         "",
		Keywords:        nil,
		ScrapeStatus:    model.UrlScraped,
	}, nil
}

func (s *ThebellScraper) parseArticle(el *rodext.Element) (*model.Article, error) {
	// todo: implement
	return nil, nil
}

// cleanThebellArticleUrl removes unnecessary query parameters from thebell article url,
// leaving only the 'key' parameter
func cleanThebellArticleUrl(u string) (string, error) {
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
