package scraper

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TheBellSuite struct {
	suite.Suite
	scraper        *TheBellScraper
	cleanupScraper func()
}

func (ts *TheBellSuite) SetupTest() {
	scraper, cleanup, err := NewTheBellScraper()
	ts.NoErrorf(err, "failed to initialize TheBellScraper: %v", err)
	ts.scraper = scraper
	ts.cleanupScraper = cleanup
}

func (ts *TheBellSuite) TearDownTest() {
	ts.cleanupScraper()
}

func (ts *TheBellSuite) TestTheBellScraper_Scrape() {
}

func TestTheBellSuite(t *testing.T) {
	suite.Run(t, new(TheBellSuite))
}

func (ts *TheBellSuite) Test_cleanTheBellArticleUrl() {
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
			got, err := cleanTheBellArticleUrl(tt.u)
			ts.Equal(tt.want, got)
			if tt.shouldError {
				ts.Error(err)
			} else {
				ts.NoError(err)
			}
		})
	}
}
