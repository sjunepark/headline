package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"net/url"
)

// ArticlesPage is a struct that contains the information of the current page of articles.
// Articles is the element that contains all the articles.
// PageNav is the element that contains the page navigation buttons.
type ArticlesPage struct {
	Keyword  string
	Articles *rodext.Element
	PageNav  *rodext.Element
	PageUrl  *url.URL
	PageNo   uint
}

func NewArticlesPage(
	keyword string,
	articles *rodext.Element,
	pageNavigation *rodext.Element,
	pageUrl *url.URL,
	pageNo uint,
) *ArticlesPage {
	return &ArticlesPage{
		Keyword:  keyword,
		Articles: articles,
		PageNav:  pageNavigation,
		PageUrl:  pageUrl,
		PageNo:   pageNo,
	}
}

func (a *ArticlesPage) Text() string {
	return a.Articles.Text()
}
