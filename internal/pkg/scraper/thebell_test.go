package scraper

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TheBellSuite struct {
	suite.Suite
	scraper *TheBellScraper
}

func (ts *TheBellSuite) SetupTest() {
	scraper, err := NewTheBellScraper()
	ts.NoErrorf(err, "failed to initialize TheBellScraper: %v", err)
	ts.scraper = scraper
}

func (ts *TheBellSuite) TearDownTest() {
	ts.scraper.Cleanup()
}

func (ts *TheBellSuite) TestTheBellScraper_Scrape() {
}

func TestTheBellSuite(t *testing.T) {
	suite.Run(t, new(TheBellSuite))
}

func (ts *TheBellSuite) Test_cleanTheBellArticleUrl() {
	tests := []struct {
		name string
		u    string
		want string
	}{
		{
			name: "happy path",
			u:    "https://thebell.co.kr/free/content/ArticleView.asp?key=202404171519391760107719&lcode=00&page=1&svccode=00",
			want: "https://thebell.co.kr/free/content/ArticleView.asp?key=202404171519391760107719",
		},
		{name: "when key is empty, it should return the original url",
			u:    "https://thebell.co.kr/free/content/ArticleView.asp?lcode=00&page=1&svccode=00",
			want: "https://thebell.co.kr/free/content/ArticleView.asp?lcode=00&page=1&svccode=00",
		},
	}
	for _, tt := range tests {
		ts.Run(tt.name, func() {
			got := cleanTheBellArticleUrl(tt.u)
			ts.Equal(tt.want, got)
		})
	}
}
