package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"net/url"
)

type ArticlesPage struct {
	PageElement *rodext.Element
	PageUrl     *url.URL
	PageNo      uint
}

func NewArticlesPage(element *rodext.Element, currentUrl *url.URL, currentPage uint) *ArticlesPage {
	return &ArticlesPage{
		PageElement: element,
		PageUrl:     currentUrl,
		PageNo:      currentPage,
	}
}

func (a *ArticlesPage) Element(selector string) (*rodext.Element, error) {
	el, err := a.PageElement.Element(selector)
	if err != nil {
		return nil, err
	}
	return el, nil
}

func (a *ArticlesPage) Text() string {
	return a.PageElement.Text()
}
