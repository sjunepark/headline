package util

import "net/url"

func IsUrlValid(u *url.URL) bool {
	return u.Scheme != "" && u.Host != ""
}
