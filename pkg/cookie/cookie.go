package cookie

import (
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// GetTargetCookie returns cookie value
// e.g. GetTargetCookie(res.Header["Set-Cookie"], "go-web-ginserver")
func GetTargetCookie(cookies []string, key string) string {
	for _, cookie := range cookies {
		tmp := strings.Split(cookie, ";")
		tmp = strings.Split(tmp[0], "=")
		if tmp[0] == key {
			return tmp[1]
		}
	}
	return ""
}

// GetTargetCookieFromString returns cookie value
func GetTargetCookieFromString(cookieURL, key string, jar *cookiejar.Jar) string {
	parsed, _ := url.Parse(cookieURL)
	cookies := jar.Cookies(parsed) // cookies []*http.Cookie

	for _, cookie := range cookies {
		if cookie.Name == key {
			return cookie.Value
		}
	}
	return ""
}
