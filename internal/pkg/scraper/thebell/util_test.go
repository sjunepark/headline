package thebell

import (
	"github.com/stretchr/testify/assert"
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
			got, err := parseDate(tc.thebellDate)
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
