package thebell

import (
	"github.com/sejunpark/headline/internal/pkg/scraper"
	"github.com/stretchr/testify/suite"
	"net/url"
	"testing"
	"time"
)

type ThebellScraperSuite struct {
	suite.Suite
	scraper *scraper.Scraper
	cleanup func()
}

func (ts *ThebellScraperSuite) SetupSuite() {
	builder, cleanup, err := NewThebellScraperBuilder()
	ts.NoErrorf(err, "failed to initialize TheBellScraperBuilder: %v", err)
	s, cleanup, err := scraper.NewScraper(builder, cleanup)
	ts.NoErrorf(err, "failed to initialize Scraper: %v", err)
	ts.scraper = s
	ts.cleanup = cleanup
}

func (ts *ThebellScraperSuite) TearDownSuite() {
	ts.cleanup()
}

func (ts *ThebellScraperSuite) TestScrape() {
	articleInfos, err := ts.scraper.Scrape("삼성전자", time.Time{})
	ts.NoErrorf(err, "failed to scrape: %v", err)
	// todo: implement assertion logic
	ts.NotEmpty(articleInfos)
}

func TestThebellScraperSuite(t *testing.T) {
	suite.Run(t, new(ThebellScraperSuite))
}

type ThebellBuilderSuite struct {
	suite.Suite
	builder *ScraperBuilder
	cleanup func()
}

func (ts *ThebellBuilderSuite) SetupSuite() {
	builder, cleanup, err := NewThebellScraperBuilder()
	ts.NoErrorf(err, "failed to initialize TheBellScraperBuilder: %v", err)
	ts.builder = builder
	ts.cleanup = cleanup
}

func (ts *ThebellBuilderSuite) TearDownSuite() {
	ts.cleanup()
}

func TestThebellBuilderSuite(t *testing.T) {
	suite.Run(t, new(ThebellBuilderSuite))
}

func (ts *ThebellBuilderSuite) Test_fetchArticlesPage() {
	ts.Run("happy path", func() {
		p, err := ts.builder.FetchArticlesPage("삼성전자", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		newsList, err := p.Element(".newsList")
		ts.NoError(err, "failed to get newsList")
		ts.NotNil(newsList)

		paging, err := p.Element(".paging")
		ts.NoError(err, "failed to get paging")
		ts.NotNil(paging)
	})

	ts.Run("when keyword is empty, it should return an error", func() {
		p, err := ts.builder.FetchArticlesPage("", time.Time{})
		ts.Error(err, "expected an error")
		ts.Nil(p, "expected nil")
	})
}

func (ts *ThebellBuilderSuite) Test_fetchNextPage() {
	ts.Run("when next page exists, it should return exists=true, and should be able to navigate to the next page, and should have different text content", func() {
		initialPage, err := ts.builder.FetchArticlesPage("삼성전자", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		nextPage, exists := ts.builder.FetchNextPage(initialPage)
		ts.True(exists)
		ts.NotNil(nextPage)

		ts.NotEqual(initialPage.Text(), nextPage.Text())
	})

	ts.Run("when there is no next page, it should return exists=false", func() {
		p, err := ts.builder.FetchArticlesPage("청호나이스", time.Time{})
		ts.NoError(err, "failed to fetch articles page")

		nextPage, exists := ts.builder.FetchNextPage(p)
		ts.False(exists)
		ts.Nil(nextPage)
	})
}

type TheBellUrlUtilSuite struct {
	suite.Suite
	urlUtil *thebellUrlUtil
}

func (ts *TheBellUrlUtilSuite) SetupTest() {
	util, err := newThebellUrlUtil()
	ts.NoErrorf(err, "failed to initialize TheBellUrlUtil: %v", err)
	ts.urlUtil = util
}

func TestTheBellUrlUtilSuite(t *testing.T) {
	suite.Run(t, new(TheBellUrlUtilSuite))
}

func (ts *TheBellUrlUtilSuite) Test_getAbsoluteUrl() {
	tests := []struct {
		name string
		u    string
		want string
	}{
		{
			name: "happy path",
			u:    "/free/content/ArticleView.asp?key=202404171519391760107719",
			want: "https://thebell.co.kr/free/content/ArticleView.asp?key=202404171519391760107719",
		},
	}
	for _, tt := range tests {
		ts.Run(tt.name, func() {
			relativeUrl, err := url.Parse(tt.u)
			ts.NoError(err, "failed to parse url")
			got := ts.urlUtil.getAbsoluteUrl(relativeUrl)
			ts.Equalf(tt.want, got.String(), "expected %s, got %s", tt.want, got)
		})
	}
}

func (ts *TheBellUrlUtilSuite) Test_getKeywordUrl() {
	tests := []struct {
		name        string
		keyword     string
		want        string
		shouldError bool
	}{
		{
			name:        "happy path when there are multiple pages",
			keyword:     "삼성전자",
			want:        "https://thebell.co.kr/free/content/Search.asp?keyword=%EC%82%BC%EC%84%B1%EC%A0%84%EC%9E%90&page=1&part=A&period=360",
			shouldError: false,
		},
		{
			name:        "happy path when there is only one page",
			keyword:     "청호나이스",
			want:        "https://thebell.co.kr/free/content/Search.asp?keyword=%EC%B2%AD%ED%98%B8%EB%82%98%EC%9D%B4%EC%8A%A4&page=1&part=A&period=360",
			shouldError: false,
		},
		{
			name:        "when keyword is empty, it should return an error",
			keyword:     "",
			want:        "",
			shouldError: true,
		},
		{
			name:        "when there are no search results, it should still return a url",
			keyword:     "하가자다아차사",
			want:        "https://thebell.co.kr/free/content/Search.asp?keyword=%ED%95%98%EA%B0%80%EC%9E%90%EB%8B%A4%EC%95%84%EC%B0%A8%EC%82%AC&page=1&part=A&period=360",
			shouldError: false,
		},
	}
	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got, err := ts.urlUtil.getKeywordUrl(tt.keyword)
			if tt.shouldError {
				ts.Nil(got, "expected nil")
				ts.Error(err, "expected an error")
			} else {
				ts.Equal(tt.want, got.String())
				ts.NoError(err, "expected no error")
			}
		})
	}

}

func (ts *TheBellUrlUtilSuite) Test_cleanArticleUrl() {
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
		ts.Run(tt.name, func() {
			got, err := ts.urlUtil.cleanArticleUrl(tt.u)
			ts.Equalf(tt.want, got, "expected %s, got %s", tt.want, got)
			if tt.shouldError {
				ts.Error(err, "expected an error")
			} else {
				ts.NoError(err, "expected no error")
			}
		})
	}
}
