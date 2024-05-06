package thebell

import (
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"github.com/sejunpark/headline/internal/pkg/scraper"
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
