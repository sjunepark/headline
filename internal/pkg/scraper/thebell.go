package scraper

import (
	"fmt"
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"net/url"
	"strings"
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

func (s *ThebellScraper) String() string {
	return "ThebellScraper"
}

type thebellUrl struct {
	baseUrl    *url.URL
	keywordUrl *url.URL
}

func newThebellUrl(keyword string) (*thebellUrl, error) {
	baseUrlString := "https://thebell.co.kr"
	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		return nil, err
	}

	keywordRelUrl := fmt.Sprintf("/free/content/Search.asp?page=1&period=360&part=A&keyword=%s", keyword)
	keywordUrl := baseUrl.ResolveReference(&url.URL{Path: keywordRelUrl})

	return &thebellUrl{
		baseUrl:    baseUrl,
		keywordUrl: keywordUrl,
	}, nil
}

func (u *thebellUrl) getAbsoluteUrl(relativeUrl *url.URL) *url.URL {
	return u.baseUrl.ResolveReference(relativeUrl)
}

type thebellArticleMetadata struct {
	title       string
	summary     string
	relativeUrl *url.URL
	datetime    time.Time
}

func (m *thebellArticleMetadata) getArticleMetadata(u *thebellUrl, keyword string) *model.ArticleMetadata {
	return &model.ArticleMetadata{
		Keywords:        map[string]bool{keyword: true},
		Title:           m.title,
		Summary:         m.summary,
		CreatedDateTime: m.datetime,
		UpdateDateTime:  m.datetime,
		Url:             u.getAbsoluteUrl(m.relativeUrl),
		Source:          "thebell",
		SourceUrl:       u.baseUrl,
	}
}

type articleElementToScrape struct {
	aTag *rodext.Element
}

func newArticleElementToScrape(el *rodext.Element) (*articleElementToScrape, error) {
	aTag, err := el.Element("a")
	if err != nil {
		return nil, err
	}
	return &articleElementToScrape{
		aTag: aTag,
	}, nil
}

func (el *articleElementToScrape) relativeUrl() (*url.URL, error) {
	href := el.aTag.Attribute("href")
	relativeUrl, err := url.Parse(href)
	if err != nil {
		return nil, err
	}
	return relativeUrl, nil
}

func (el *articleElementToScrape) title() (string, error) {
	title := el.aTag.Attribute("title")
	if title == "" {
		return "", fmt.Errorf("title attribute not found")
	}
	return title, nil
}

func (el *articleElementToScrape) summary() (string, error) {
	summaryEl, err := el.aTag.Element("dd")
	if err != nil {
		return "", err
	}
	return summaryEl.Text(), nil
}

func (el *articleElementToScrape) datetime() (time.Time, error) {
	datetimeEl, err := el.aTag.Element("dd>.date")
	if err != nil {
		return time.Time{}, err
	}
	koreanDatetime := datetimeEl.Text()

	// Replace Korean AM/PM with standard AM/PM
	replacements := map[string]string{
		"오전": "AM",
		"오후": "PM",
	}
	for k, v := range replacements {
		koreanDatetime = strings.Replace(koreanDatetime, k, v, 1)
	}

	// Define the layout and parse the time
	// Note: "2006-01-02 3:04:05 PM" is the reference time format used by Go
	const layout = "2006-01-02 3:04:05 PM"
	parsedTime, err := time.Parse(layout, koreanDatetime)
	if err != nil {
		return time.Time{}, err
	}

	// Set timezone, assuming you want to convert it to KST (Korea Standard Time)
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.Time{}, err
	}
	kstTime := parsedTime.In(location)

	return kstTime, nil
}

func (el *articleElementToScrape) articleMetadata(u *thebellUrl, keyword string) (*model.ArticleMetadata, error) {
	relativeUrl, err := el.relativeUrl()
	if err != nil {
		return nil, err
	}
	title, err := el.title()
	if err != nil {
		return nil, err
	}
	summary, err := el.summary()
	if err != nil {
		return nil, err
	}
	datetime, err := el.datetime()
	if err != nil {
		return nil, err
	}

	return &model.ArticleMetadata{
		Keywords:        map[string]bool{keyword: true},
		Title:           title,
		Summary:         summary,
		CreatedDateTime: datetime,
		UpdateDateTime:  datetime,
		Url:             u.getAbsoluteUrl(relativeUrl),
		Source:          "thebell",
		SourceUrl:       u.baseUrl,
	}, nil
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
