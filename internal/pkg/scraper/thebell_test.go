package scraper

import (
	"fmt"
	"github.com/sejunpark/headline/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ScraperSuite struct {
	suite.Suite
	keywords []string
	scrapers []*Scraper
	cleanups []func()
}

func (ts *ScraperSuite) SetupTest() {
	ts.keywords = []string{"삼성전자", "SK하이닉스"}

	// cleanup will be handled by the Scraper
	thebellSourceScraper, _, err := NewThebellScraper()
	ts.NoErrorf(err, "failed to initialize thebellSourceScraper: %v", err)
	thebellScraper, thebellCleanup := NewScraper(thebellSourceScraper)
	ts.scrapers = append(ts.scrapers, thebellScraper)
	ts.cleanups = append(ts.cleanups, thebellCleanup)
}

func (ts *ScraperSuite) TearDownTest() {
	for _, cleanup := range ts.cleanups {
		cleanup()
	}
}

func TestScraperSuite(t *testing.T) {
	suite.Run(t, new(ScraperSuite))
}

func (ts *ScraperSuite) TestScraper_fetchArticles() {
	for _, scraper := range ts.scrapers {
		ts.Run(fmt.Sprintf("fetchArticles for %s", scraper), func() {
			var articles chan *model.ArticleMetadata

			for _, keyword := range ts.keywords {
				go func(keyword string) {
					articlesToSend, err := scraper.fetchArticles(keyword, time.Time{})
					ts.NoErrorf(err, "failed to fetch articles for keyword %s: %v", keyword, err)
					for article := range articlesToSend {
						articles <- article
					}
				}(keyword)
			}

			for article := range articles {
				ts.Truef(article.IsValid(), "invalid article metadata: %v", article)
			}
		})
	}
}

func Test_cleanThebellUrl(t *testing.T) {
	tests := []struct {
		name        string
		u           string
		want        string
		shouldError bool
	}{
		{
			name:        "happy path",
			u:           "https://scraper.co.kr/free/content/ArticleView.asp?key=202404171519391760107719&lcode=00&page=1&svccode=00",
			want:        "https://scraper.co.kr/free/content/ArticleView.asp?key=202404171519391760107719",
			shouldError: false,
		},
		{name: "when there is no key parameter, it should return an error",
			u:           "https://scraper.co.kr/free/content/ArticleView.asp?lcode=00&page=1&svccode=00",
			want:        "",
			shouldError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cleanThebellUrl(tt.u)
			assert.Equal(t, tt.want, got)
			if tt.shouldError {
				t.Error(err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
