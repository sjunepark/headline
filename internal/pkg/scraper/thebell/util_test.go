package thebell

import (
	"fmt"
	"github.com/sejunpark/headline/internal/pkg/rodext"
	"github.com/sejunpark/headline/internal/pkg/scraper"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func Test_parseDate(t *testing.T) {
	kst, err := time.LoadLocation("Asia/Seoul")
	assert.NoError(t, err)

	tt := []struct {
		name        string
		thebellDate string
		want        time.Time
		shouldError bool
	}{
		{
			name:        "happy path",
			thebellDate: "2023-10-04 오전 7:34:13",
			want:        time.Date(2023, 10, 4, 7, 34, 13, 0, kst),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseDatetime(tc.thebellDate)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Equal(t, time.Time{}, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

}

func Test_currentPageNoIsValid(t *testing.T) {
	browser, cleanup, err := rodext.NewBrowser(rodext.DefaultBrowserOptions)
	assert.NoError(t, err)
	defer cleanup()

	page, _, err := browser.Page()
	assert.NoError(t, err)

	tt := []struct {
		name    string
		keyword string
		pageNo  uint
		want    bool
	}{
		{
			name:    "should return true when page number is valid",
			keyword: "삼성전자",
			pageNo:  1,
			want:    true,
		},
		{
			name:    "should return false when page number is not valid",
			keyword: "nonexistingcompany",
			pageNo:  100,
			want:    false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			urlString := fmt.Sprintf("https://thebell.co.kr/free/content/Search.asp?page=%d&sdt=&period=360&part=A&keyword=%s", tc.pageNo, tc.keyword)
			wait, navErr := page.Navigate(urlString)
			assert.NoError(t, navErr)
			waitErr := wait(".newsBox")
			assert.NoError(t, waitErr)

			el, elErr := page.Element(".newsBox")
			assert.NoError(t, elErr)

			u, parseErr := url.Parse(urlString)
			assert.NoError(t, parseErr)

			articlesPage := scraper.NewArticlesPage(tc.keyword, el, u, tc.pageNo)
			assert.True(t, currentPageNoIsValid(articlesPage) == tc.want)
		})
	}
}
