package rodext

import (
	"fmt"
	"net/url"
)

type thebellUrlUtil struct {
	base    *url.URL
	keyword *url.URL
}

func newThebellUrlUtil(keyword string) (*thebellUrlUtil, error) {
	baseUrlString := "https://thebell.co.kr"
	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		return nil, err
	}

	keywordRelUrl := fmt.Sprintf("/free/content/Search.asp?page=1&period=360&part=A&keyword=%s", keyword)
	keywordUrl := baseUrl.ResolveReference(&url.URL{Path: keywordRelUrl})

	return &thebellUrlUtil{
		base:    baseUrl,
		keyword: keywordUrl,
	}, nil
}

func (t *thebellUrlUtil) getKeywordUrl(keyword string) (*url.URL, error) {
	unparsedUrl := fmt.Sprintf("https://thebell.co.kr/free/content/Search.asp?page=1&period=360&part=A&keyword=%s", keyword)
	parsedUrl, err := url.Parse(unparsedUrl)
	if err != nil {
		return nil, err
	}
	return parsedUrl, nil
}

func (t *thebellUrlUtil) cleanArticleUrl(u *url.URL) (*url.URL, error) {
	query := u.Query()
	key := query.Get("key")
	if key == "" {
		return nil, fmt.Errorf("parameter 'key' not found in url: %s", u)
	}
	query = url.Values{"key": []string{key}}
	u.RawQuery = query.Encode()
	return u, nil
}
