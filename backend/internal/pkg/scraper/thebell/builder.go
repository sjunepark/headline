package thebell

import (
	"github.com/cockroachdb/errors"
	"github.com/sejunpark/headline/backend/internal/pkg/model"
	"github.com/sejunpark/headline/backend/internal/pkg/rodext"
	"log/slog"
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

	browser, cleanup, err := rodext.NewBrowser(rodext.DefaultBrowserOptions)
	if err != nil {
		return nil, nil, err
	}

	return &ScraperBuilder{browser: browser, util: util}, cleanup, nil
}

func (b *ScraperBuilder) FetchArticlesPage(keyword string, _ time.Time) (*model.ArticlesPage, error) {
	functionName := "FetchArticlesPage"

	keywordUrl, err := b.util.getKeywordUrl(keyword)
	if err != nil {
		return nil, errors.Wrapf(err, "getKeywordUrl(%s) failed", keyword)
	}
	pageNo, err := b.util.getPageNo(keywordUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "getPageNo(%s) failed", keywordUrl.String())
	}

	page, _, err := b.browser.Page()
	if err != nil {
		return nil, errors.Wrap(err, "browser.Page() failed")
	}
	wait, err := page.Navigate(keywordUrl.String())
	if err != nil {
		return nil, errors.Wrapf(err, "page.Navigate(%s) failed", keywordUrl.String())
	}

	articlesEl, pageNavEl, err := getPageElements(page, wait)
	if err != nil {
		return nil, errors.Wrap(err, "getPageElements() failed")
	}

	articlesPage := model.NewArticlesPage(keyword, articlesEl, pageNavEl, keywordUrl, pageNo)
	slog.Debug("executed:", "function", functionName, "keyword", keyword, "pageNo", pageNo)
	return articlesPage, nil
}

func (b *ScraperBuilder) FetchNextPage(currentPage *model.ArticlesPage) (nextPage *model.ArticlesPage, err error) {
	functionName := "FetchNextPage"

	if currentPage == nil {
		return nil, errors.AssertionFailedf("currentPage is nil")
	}

	nextPageUrl, nextPageNo, err := b.util.getNextPageUrl(currentPage.PageUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "getNextPageUrl(%s) failed", currentPage.PageUrl.String())
	}

	p, _, err := b.browser.Page()
	if err != nil {
		return nil, errors.Wrap(err, "browser.Page() failed")
	}
	wait, err := p.Navigate(nextPageUrl.String())
	if err != nil {
		return nil, errors.Wrapf(err, "page.Navigate(%s) failed", nextPageUrl.String())
	}
	slog.Debug("navigated to next page.", "function", functionName, "keyword", currentPage.Keyword, "pageNo", nextPageNo)

	nextPageEl, pageNavEl, err := getPageElements(p, wait)
	if err != nil {
		return nil, errors.Wrap(err, "getPageElements() failed")
	}

	nextPage = model.NewArticlesPage(currentPage.Keyword, nextPageEl, pageNavEl, nextPageUrl, nextPageNo)
	if !currentPageNoIsValid(nextPage) {
		return nil, errors.Newf("currentPageNoIsValid(%v) failed", nextPage)
	}

	return nextPage, nil
}

func getPageElements(page *rodext.Page, wait func(selector string) error) (articlesEl *rodext.Element, pageNavEl *rodext.Element, err error) {
	articlesSelector := ".listBox>ul"
	err = wait(articlesSelector)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "wait(%s) failed", articlesSelector)
	}
	articlesEl, err = page.Element(articlesSelector)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "page.Element(%s) failed", articlesSelector)
	}

	pageNavSelector := ".paging"
	err = wait(pageNavSelector)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "wait(%s) failed", pageNavSelector)
	}
	pageNavEl, err = page.Element(pageNavSelector)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "page.Element(%s) failed", pageNavSelector)
	}

	return articlesEl, pageNavEl, nil
}

func (b *ScraperBuilder) ParseArticlesPage(p *model.ArticlesPage) ([]*model.ArticleInfo, error) {
	functionName := "ParseArticlesPage"

	dlTags, err := p.Articles.Elements("li>dl")
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
			slog.Error("failed to get aTag.", "function", functionName, "error", elErr)
			continue
		}

		// articleInfo.Title
		title, attrErr := aTag.Attribute("title")
		if attrErr != nil {
			slog.Error("failed to get title attribute.", "function", functionName, "error", attrErr)
			continue
		}
		if title == "" {
			slog.Error("title is empty.", "function", functionName)
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
			slog.Error("failed to get dateTag.", "function", functionName, "error", elErr)
			continue
		}
		datetime, dateErr := parseDatetime(dateTag.Text())
		if dateErr != nil {
			slog.Error("failed to parse datetime.", "function", functionName, "error", dateErr)
			continue
		}
		articleInfo.CreatedDateTime = datetime
		articleInfo.UpdateDateTime = datetime

		// articleInfo.Url
		relUrlStr, elErr := aTag.Attribute("href")
		if elErr != nil || relUrlStr == "" {
			slog.Error("failed to get href attribute.", "function", functionName, "error", elErr)
			continue
		}
		absUrl, urlErr := b.util.getAbsoluteUrl(relUrlStr)
		if urlErr != nil {
			slog.Error("failed to get absolute url.", "function", functionName, "error", urlErr)
			continue
		}
		articleInfo.Url = absUrl

		articleInfos = append(articleInfos, articleInfo)
		slog.Debug("appended ArticleInfo.", "function", functionName, "title", articleInfo.Title)

	}
	return articleInfos, nil
}
