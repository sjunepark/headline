package scraper

import (
	"fmt"
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"log/slog"
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

func (s *ThebellScraper) fetchArticles(keyword string, startDate time.Time) (<-chan *model.ArticleMetadata, error) {
	page, putPage, err := s.browser.Page()
	if err != nil {
		return nil, err
	}
	defer putPage()

	fullUrl := fmt.Sprintf("https://thebell.co.kr/free/content/Search.asp?page=1&period=360&part=A&keyword=%s", keyword)
	err = page.Navigate(fullUrl)
	if err != nil {
		return nil, err
	}

	parsedFullUrl, err := url.Parse(fullUrl)
	if err != nil {
		return nil, err
	}
	sourceUrl := &url.URL{Host: parsedFullUrl.Host, Scheme: parsedFullUrl.Scheme}

	articleElementsToScrape, err := page.Elements("div.listBox>ul>li>dl")
	if err != nil {
		return nil, err
	}

	articles := make(chan *model.ArticleMetadata)
	go func() {
		defer close(articles)
		for _, articleElementToScrape := range articleElementsToScrape {
			articleMetadata, parseErr := parseArticleElementToScrape(articleElementToScrape)
			if parseErr != nil {
				slog.Error("failed to parse article element", "error", parseErr)
				continue
			}

			articles <- &model.ArticleMetadata{
				Keywords:        map[string]bool{keyword: true},
				Title:           articleMetadata.title,
				Summary:         articleMetadata.summary,
				CreatedDateTime: articleMetadata.datetime,
				UpdateDateTime:  articleMetadata.datetime,
				Url:             sourceUrl.ResolveReference(articleMetadata.relativeUrl),
				Source:          "thebell",
				SourceUrl:       sourceUrl,
			}
			slog.Info("scraped and sent article to channel", "title", articleMetadata.title)
		}
	}()
	return articles, nil
}

// cleanThebellUrl removes unnecessary query parameters from thebell article url,
// leaving only the 'key' parameter
func cleanThebellUrl(u string) (string, error) {
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

func (s *ThebellScraper) String() string {
	return "ThebellScraper"
}

type thebellArticleMetadata struct {
	title       string
	summary     string
	relativeUrl *url.URL
	datetime    time.Time
}

func parseArticleElementToScrape(el *rodext.Element) (*thebellArticleMetadata, error) {
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

	summaryEl, summaryParesErr := a.Element("dd")
	if summaryParesErr != nil {
		return nil, summaryParesErr
	}
	summary := summaryEl.Text()

	datetimeEl, datetimeElParseErr := el.Element("dd>.date")
	if datetimeElParseErr != nil {
		return nil, datetimeElParseErr
	}

	datetime, err := parseThebellKoreanDatetime(datetimeEl.Text())
	if err != nil {
		return nil, err
	}

	return &thebellArticleMetadata{
		title:       title,
		summary:     summary,
		relativeUrl: relativeUrl,
		datetime:    datetime,
	}, nil
}
