package scraper

import (
	"strings"
	"time"
)

// parseThebellKoreanDatetime parses a Korean datetime string in format "2023-06-28 오후 12:58:46" to a time.Time object.
func parseThebellKoreanDatetime(koreanDatetime string) (time.Time, error) {
	// Replace Korean AM/PM with standard AM/PM
	replacements := map[string]string{
		"오전": "AM",
		"오후": "PM",
	}
	for k, v := range replacements {
		koreanDatetime = strings.Replace(koreanDatetime, k, v, 1)
	}

	// Define the layout and parse the time
	// Note: "2006-01-02 3:04:05 PM" is the reference time format used by Go
	const layout = "2006-01-02 3:04:05 PM"
	parsedTime, err := time.Parse(layout, koreanDatetime)
	if err != nil {
		return time.Time{}, err
	}

	// Set timezone, assuming you want to convert it to KST (Korea Standard Time)
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.Time{}, err
	}
	kstTime := parsedTime.In(location)

	return kstTime, nil
}
