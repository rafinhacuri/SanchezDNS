package utils

import (
	"net/url"
	"strings"
)

func NormalizeBase(h string) string {
	b := strings.TrimRight(h, "/")
	if !strings.HasPrefix(b, "http://") && !strings.HasPrefix(b, "https://") {
		b = "http://" + b
	}
	u, err := url.Parse(b)
	if err != nil {
		return b
	}
	if u.Port() == "" {
		u.Host = u.Host + ":8081"
	}
	return u.String()
}
