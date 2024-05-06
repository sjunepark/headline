package thebell

import (
	"fmt"
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"github.com/sejunpark/headline/internal/pkg/scraper"
	"log/slog"
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
	ap := scraper.NewArticlesPage(keyword, el, keywordUrl, pageNo)
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

	return scraper.NewArticlesPage(currentPage.Keyword, nextPageEl, nextPageUrl, currentPage.PageNo+1), true
}

func (b *ScraperBuilder) ParseArticlesPage(p *scraper.ArticlesPage) ([]*model.ArticleInfo, error) {
	dlTags, err := p.Elements("ul>li>dl")
	if err != nil {
		return nil, err
	}

	var articleInfos []*model.ArticleInfo
	for _, dlTag := range dlTags {
		articleInfo := &model.ArticleInfo{
			Keywords:        map[string]bool{p.Keyword: true},
			Title:           "",
			Summary:         "",
			CreatedDateTime: time.Time{},
			UpdateDateTime:  time.Time{},
			Url:             nil,
			Source:          "thebell",
			SourceUrl:       b.util.base,
		}

		aTag, elErr := dlTag.Element("a")
		if elErr != nil {
			slog.Error("failed to get aTag", "error", elErr)
			continue
		}

		// articleInfo.Title
		title, attrErr := aTag.Attribute("title")
		if attrErr != nil {
			slog.Error("failed to get title attribute", "error", attrErr)
			continue
		}
		if title == "" {
			slog.Error("title is empty")
			continue
		}
		articleInfo.Title = title

		// articleInfo.Summary
		summaryTag, elErr := aTag.Element("dd")
		var summary string
		if elErr != nil {
			summary = ""
		} else {
			summary = summaryTag.Text()
		}
		articleInfo.Summary = summary

		// articleInfo.CreatedDateTime, articleInfo.UpdateDateTime
		dateTag, elErr := dlTag.Element(".date")
		if elErr != nil {
			slog.Error("failed to get dateTag", "error", elErr)
			continue
		}
		datetime, dateErr := parseDatetime(dateTag.Text())
		if dateErr != nil {
			slog.Error("failed to parse datetime", "error", dateErr)
			return nil, dateErr
		}
		articleInfo.CreatedDateTime = datetime
		articleInfo.UpdateDateTime = datetime

		// articleInfo.Url
		relUrlStr, elErr := aTag.Attribute("href")
		if elErr != nil || relUrlStr == "" {
			slog.Error("failed to get href attribute", "error", elErr)
			continue
		}
		absUrl, urlErr := b.util.getAbsoluteUrl(relUrlStr)
		if urlErr != nil {
			slog.Error("failed to get absolute url", "error", urlErr)
			continue
		}
		articleInfo.Url = absUrl

		articleInfos = append(articleInfos, articleInfo)
	}
	return articleInfos, nil
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

// getAbsoluteUrl returns an error if it fails to parse the relativeUrl
func (util *thebellUrlUtil) getAbsoluteUrl(relativeUrl string) (*url.URL, error) {
	u, err := url.Parse(relativeUrl)
	if err != nil {
		return nil, err
	}
	return util.base.ResolveReference(u), nil
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
func (util *thebellUrlUtil) cleanArticleUrl(articleUrl string) (string, error) {
	parsedUrl, err := url.Parse(articleUrl)
	if err != nil {
		return "", err
	}
	query := parsedUrl.Query()
	key := query.Get("key")
	if key == "" {
		return "", fmt.Errorf("parameter 'key' not found in url: %s", articleUrl)
	}
	query = url.Values{"key": []string{key}}
	parsedUrl.RawQuery = query.Encode()
	articleUrl = parsedUrl.String()
	return articleUrl, nil
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
