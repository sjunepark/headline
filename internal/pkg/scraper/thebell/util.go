package thebell

import (
	"strings"
	"time"
)

// parseDatetime parses the date string from thebell.co.kr
// thebellDate format: "2023-10-04 오전 7:34:13"
func parseDatetime(thebellDate string) (time.Time, error) {
	translated := strings.Replace(thebellDate, "오전", "AM", 1)
	translated = strings.Replace(translated, "오후", "PM", 1)

	kst, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.Time{}, err
	}

	layout := "2006-01-02 PM 3:04:05"

	parsed, err := time.ParseInLocation(layout, translated, kst)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}
