package scraper

import (
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"net/url"
)

type ArticlesPage struct {
	Keyword     string
	PageElement *rodext.Element
	PageUrl     *url.URL
	PageNo      uint
}

func NewArticlesPage(keyword string, element *rodext.Element, currentUrl *url.URL, currentPageNo uint) *ArticlesPage {
	return &ArticlesPage{
		Keyword:     keyword,
		PageElement: element,
		PageUrl:     currentUrl,
		PageNo:      currentPageNo,
	}
}

func (a *ArticlesPage) Element(selector string) (*rodext.Element, error) {
	el, err := a.PageElement.Element(selector)
	if err != nil {
		return nil, err
	}
	return el, nil
}

func (a *ArticlesPage) Elements(selector string) ([]*rodext.Element, error) {
	els, err := a.PageElement.Elements(selector)
	if err != nil {
		return nil, err
	}
	return els, nil
}

func (a *ArticlesPage) Text() string {
	return a.PageElement.Text()
}
