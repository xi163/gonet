package uri

import (
	"strings"
)

func ParseURL(s string, v any) {
	if strings.Contains(s, "?") {
		ParseQuery(s[strings.Index(s, "?")+1:], v)
	} else {
		if strings.Contains(s, "/") {
		} else if strings.Contains(s, "=") {
			ParseQuery(s, v)
		}
	}
}

func ParseQuery(query string, v any) {
	// values, _ := url.ParseQuery(query)
	// for url := range values {
	// 	sub := url[strings.Index(url, "?")+1:]
	// 	dic := map[string]string{}
	// 	for {
	// 		s := strings.Index(sub, "=")
	// 		if s == -1 {
	// 			break
	// 		}
	// 		p := strings.Index(sub, "&")
	// 		if p == -1 {
	// 			dic[sub[0:s]] = sub[s+1:]
	// 			break
	// 		} else {
	// 			dic[sub[0:s]] = sub[s+1 : p]
	// 		}
	// 		sub = sub[p+1:]
	// 	}
	//
	//
	// url := query
	// sub := url[strings.Index(url, "?")+1:]
	// dic := map[string]string{}
	// for {
	// 	s := strings.Index(sub, "=")
	// 	if s == -1 {
	// 		break
	// 	}
	// 	p := strings.Index(sub, "&")
	// 	if p == -1 {
	// 		dic[sub[0:s]] = sub[s+1:]
	// 		break
	// 	} else if p > s {
	// 		dic[sub[0:s]] = sub[s+1 : p]
	// 		sub = sub[p+1:]
	// 	} else {
	// 		break
	// 	}
	// }
	switch v := v.(type) {
	case *map[string]any:
		for {
			s := strings.Index(query, "=")
			if s == -1 {
				break
			}
			p := strings.Index(query, "&")
			if p == -1 {
				(*v)[query[0:s]] = query[s+1:]
				break
			} else if p > s {
				(*v)[query[0:s]] = query[s+1 : p]
				query = query[p+1:]
			} else {
				break
			}
		}
	}
}
