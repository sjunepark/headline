package scraper

import (
	"fmt"
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"net/url"
	"time"
)

type thebellScraperBuilder struct {
	browser *rodext.Browser
	util    *thebellUrlUtil
}

func newThebellScraperBuilder() (builder *thebellScraperBuilder, cleanup func(), err error) {
	util, err := newThebellUrlUtil()
	if err != nil {
		return nil, nil, err
	}

	browser, cleanup, err := rodext.NewBrowser(browserOptions)
	if err != nil {
		return nil, nil, err
	}

	return &thebellScraperBuilder{browser: browser, util: util}, cleanup, nil
}

func (b *thebellScraperBuilder) fetchArticlesPage(keyword string, startDate time.Time) (*rodext.Element, error) {
	keywordUrl, err := b.util.getKeywordUrl(keyword)
	if err != nil {
		return nil, err
	}
	page, _, err := b.browser.Page()
	if err != nil {
		return nil, err
	}
	err = page.Navigate(keywordUrl.String())
	if err != nil {
		return nil, err
	}
	el, err := page.Element(".newsBox")
	if err != nil {
		return nil, err
	}
	return el, nil
}

func (b *thebellScraperBuilder) fetchNextPage(currentPage *rodext.Element) (p *rodext.Element, exists bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *thebellScraperBuilder) parseArticlesPage(p *rodext.Element) ([]*model.ArticleInfos, error) {
	//TODO implement me
	panic("implement me")
}

type thebellUrlUtil struct {
	base *url.URL
}

func newThebellUrlUtil() (*thebellUrlUtil, error) {
	baseUrlString := "https://thebell.co.kr"
	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		return nil, err
	}
	return &thebellUrlUtil{base: baseUrl}, nil
}

func (util *thebellUrlUtil) getAbsoluteUrl(relativeUrl *url.URL) *url.URL {
	return util.base.ResolveReference(relativeUrl)
}

func (util *thebellUrlUtil) getKeywordUrl(keyword string) (*url.URL, error) {
	if keyword == "" {
		return nil, fmt.Errorf("keyword is empty")
	}
	keywordRelPath := "/free/content/Search.asp"

	params := url.Values{}
	params.Set("page", "1")
	params.Set("period", "360")
	params.Set("part", "A")
	params.Set("keyword", keyword)

	keywordUrl := &url.URL{
		Path:     keywordRelPath,
		RawQuery: params.Encode(),
	}

	absoluteUrl := util.base.ResolveReference(keywordUrl)
	return absoluteUrl, nil
}

// cleanArticleUrl removes unnecessary query parameters from thebell article url,
// leaving only the 'key' parameter
func (util *thebellUrlUtil) cleanArticleUrl(u string) (string, error) {
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
