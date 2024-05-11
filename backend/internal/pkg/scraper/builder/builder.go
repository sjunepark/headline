package builder

import (
	"github.com/sejunpark/headline/backend/internal/pkg/model"
	"time"
)

// Builder is an interface that defines the methods which a scraper builder must implement.
//
// FetchArticlesPage should return an ArticlesPage object for the given keyword and start date.
// The implementation of ArticlesPage is up to the builder, and effects how other methods handle it.
//
// FetchNextPage should return the next ArticlesPage if it exists.
// It should also check if the current page and the next page are different.
// If not, return false.
// Since an error is not returned from this function, it should log all errors.
//
// ParseArticlesPage should parse and return model.ArticleInfo objects.
type Builder interface {
	FetchArticlesPage(keyword string, startDate time.Time) (*model.ArticlesPage, error)
	FetchNextPage(currentPage *model.ArticlesPage) (nextPage *model.ArticlesPage, err error)
	ParseArticlesPage(p *model.ArticlesPage) ([]*model.ArticleInfo, error)
}
