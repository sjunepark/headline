package thebell

import (
	"github.com/sejunpark/headline/backend/internal/pkg/model"
	"strconv"
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

// currentPageNoIsValid checks if the current page number is valid.
// This check is used since the bell shows the last page when it's queried with a large page number beyond range.
func currentPageNoIsValid(p *model.ArticlesPage) bool {
	pageNav := p.PageNav

	// When properly accessed, thebell colors the current page number element,
	// which is represented as an "em.cur" element.
	currentPageEl, err := pageNav.Element("em.cur")
	if err == nil {
		currentPageNo, convErr := strconv.Atoi(currentPageEl.Text())
		if convErr == nil && currentPageNo == int(p.PageNo) {
			return true
		}
	}

	// If em.cur is not found, check the page numbers in the page navigation anchor elements.
	elements, err := pageNav.Elements(".paging>.btnPage")
	if err != nil {
		return false
	}
	for _, el := range elements {
		attribute, attrErr := el.Attribute("id")
		if attrErr != nil {
			return false
		}
		no, convErr := strconv.Atoi(attribute)
		if convErr != nil {
			return false
		}

		if p.PageNo == uint(no) {
			return true
		}
	}
	return false
}
