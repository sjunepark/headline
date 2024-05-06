package thebell

import (
	"fmt"
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"github.com/sejunpark/headline/internal/pkg/scraper"
	"net/url"
	"strconv"
	"time"
)

type ScraperBuilder struct {
	browser *rodext.Browser
	util    *thebellUrlUtil
}

func NewThebellScraperBuilder() (builder *ScraperBuilder, cleanup func(), err error) {
	util, err := newThebellUrlUtil()
	if err != nil {
		return nil, nil, err
	}

	browser, cleanup, err := rodext.NewBrowser(scraper.DefaultBrowserOptions)
	if err != nil {
		return nil, nil, err
	}

	return &ScraperBuilder{browser: browser, util: util}, cleanup, nil
}

func (b *ScraperBuilder) FetchArticlesPage(keyword string, startDate time.Time) (*scraper.ArticlesPage, error) {
	keywordUrl, err := b.util.getKeywordUrl(keyword)
	if err != nil {
		return nil, err
	}
	pageNo, err := b.util.getPageNo(keywordUrl)
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
	ap := scraper.NewArticlesPage(el, keywordUrl, pageNo)
	return ap, nil
}

func (b *ScraperBuilder) FetchNextPage(currentPage *scraper.ArticlesPage) (nextPage *scraper.ArticlesPage, exists bool) {
	nextPageUrl, err := b.util.getNextPageUrl(currentPage.PageUrl)
	if err != nil {
		return nil, false
	}

	p, _, err := b.browser.Page()
	if err != nil {
		return nil, false
	}
	err = p.Navigate(nextPageUrl.String())
	if err != nil {
		return nil, false
	}

	nextPageEl, err := p.Element(".newsBox")
	if err != nil {
		return nil, false
	}
	if equal, equalErr := nextPageEl.Equal(currentPage.PageElement); equalErr != nil || equal {
		return nil, false
	}

	return scraper.NewArticlesPage(nextPageEl, nextPageUrl, currentPage.PageNo+1), true
}

func (b *ScraperBuilder) ParseArticlesPage(p *scraper.ArticlesPage) ([]*model.ArticleInfos, error) {
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

func (util *thebellUrlUtil) getPageNo(u *url.URL) (uint, error) {
	query := u.Query()
	pageNoStr := query.Get("page")
	if pageNoStr == "" {
		return 1, nil
	}
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		return 0, err
	}
	return uint(pageNo), nil
}

func (util *thebellUrlUtil) getNextPageUrl(currentUrl *url.URL) (*url.URL, error) {
	pageNo, err := util.getPageNo(currentUrl)
	if err != nil {
		return nil, err
	}
	nextPageNo := pageNo + 1

	query := currentUrl.Query()
	query.Set("page", strconv.Itoa(int(nextPageNo)))
	currentUrl.RawQuery = query.Encode()
	return currentUrl, nil
}
