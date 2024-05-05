package thebell

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			got, err := cleanThebellArticleUrl(tt.u)
			assert.Equal(t, tt.want, got)
			if tt.shouldError {
				assert.Errorf(t, err, "expected an error")
			} else {
				assert.NoErrorf(t, err, "expected no error")
			}
		})
	}
}
