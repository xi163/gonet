package uri

import "net/url"

func URLEncode(s string) string {
	// uri, _ := url.Parse(s)
	// return uri.EscapedPath()
	return url.QueryEscape(s)
}

func URLDecode(s string) string {
	d, _ := url.QueryUnescape(s)
	return d
}
