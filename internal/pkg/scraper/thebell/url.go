package thebell

import (
	"github.com/cockroachdb/errors"
	"net/url"
	"strconv"
)

type thebellUrlUtil struct {
	base *url.URL
}

func newThebellUrlUtil() (*thebellUrlUtil, error) {
	baseUrlString := "https://thebell.co.kr"
	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		return nil, err
	}
	return &thebellUrlUtil{base: baseUrl}, nil
}

// getAbsoluteUrl returns an error if it fails to parse the relativeUrl
func (util *thebellUrlUtil) getAbsoluteUrl(relativeUrl string) (*url.URL, error) {
	u, err := url.Parse(relativeUrl)
	if err != nil {
		return nil, err
	}
	return util.base.ResolveReference(u), nil
}

func (util *thebellUrlUtil) getKeywordUrl(keyword string) (*url.URL, error) {
	if keyword == "" {
		return nil, errors.New("keyword is empty")
	}
	keywordRelPath := "/free/content/Search.asp"

	params := url.Values{}
	params.Set("page", "1")
	params.Set("period", "360")
	params.Set("part", "A")
	params.Set("keyword", keyword)

	keywordUrl := &url.URL{
		Path:     keywordRelPath,
		RawQuery: params.Encode(),
	}

	absoluteUrl := util.base.ResolveReference(keywordUrl)
	return absoluteUrl, nil
}

// cleanArticleUrl removes unnecessary query parameters from thebell article url,
// leaving only the 'key' parameter
func (util *thebellUrlUtil) cleanArticleUrl(articleUrl string) (string, error) {
	parsedUrl, err := url.Parse(articleUrl)
	if err != nil {
		return "", err
	}
	query := parsedUrl.Query()
	key := query.Get("key")
	if key == "" {
		return "", errors.Newf("parameter 'key' not found in url: %s", articleUrl)
	}
	query = url.Values{"key": []string{key}}
	parsedUrl.RawQuery = query.Encode()
	articleUrl = parsedUrl.String()
	return articleUrl, nil
}

func (util *thebellUrlUtil) getPageNo(u *url.URL) (uint, error) {
	query := u.Query()
	pageNoStr := query.Get("page")
	if pageNoStr == "" {
		return 1, nil
	}
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil {
		return 0, err
	}
	return uint(pageNo), nil
}

func (util *thebellUrlUtil) getNextPageUrl(currentUrl *url.URL) (u *url.URL, nextPageNo uint, err error) {
	pageNo, err := util.getPageNo(currentUrl)
	if err != nil {
		return nil, 0, err
	}
	nextPageNo = pageNo + 1

	query := currentUrl.Query()
	query.Set("page", strconv.Itoa(int(nextPageNo)))
	currentUrl.RawQuery = query.Encode()
	return currentUrl, nextPageNo, nil
}
