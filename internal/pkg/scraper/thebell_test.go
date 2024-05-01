package scraper

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ThebellSuite struct {
	suite.Suite
	scraper        *ThebellScraper
	cleanupScraper func()
}

func (ts *ThebellSuite) SetupTest() {
	scraper, cleanup, err := NewThebellScraper()
	ts.NoErrorf(err, "failed to initialize ThebellScraper: %v", err)
	ts.scraper = scraper
	ts.cleanupScraper = cleanup
}

func (ts *ThebellSuite) TearDownTest() {
	ts.cleanupScraper()
}

func TestThebellSuite(t *testing.T) {
	suite.Run(t, new(ThebellSuite))
}

func (ts *ThebellSuite) Test_cleanThebellArticleUrl() {
	tests := []struct {
		name        string
		u           string
		want        string
		shouldError bool
	}{
		{
			name:        "happy path",
			u:           "https://thebell.co.kr/free/content/ArticleView.asp?key=202404171519391760107719&lcode=00&page=1&svccode=00",
			want:        "https://thebell.co.kr/free/content/ArticleView.asp?key=202404171519391760107719",
			shouldError: false,
		},
		{name: "when there is no key parameter, it should return an error",
			u:           "https://thebell.co.kr/free/content/ArticleView.asp?lcode=00&page=1&svccode=00",
			want:        "",
			shouldError: true,
		},
	}
	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got, err := cleanThebellArticleUrl(tt.u)
			ts.Equal(tt.want, got)
			if tt.shouldError {
				ts.Error(err)
			} else {
				ts.NoError(err)
			}
		})
	}
}
