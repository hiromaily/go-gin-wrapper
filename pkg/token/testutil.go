package token

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// GetToken returns token string from parsed html for test use
func GetToken(client *http.Client, urlString string) (string, error) {
	req, _ := http.NewRequest("GET", urlString, nil)
	res, err := client.Do(req)
	defer func() {
		if res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return "", err
	}
	return parseToken(res)
}

func parseToken(res *http.Response) (string, error) {
	//<input type="hidden" name="gintoken" value="{{ .gintoken }}">
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var token string
	doc.Find("input[name=gintoken]").Each(func(_ int, s *goquery.Selection) {
		if val, ok := s.Attr("value"); ok {
			token = val
		}
	})
	return token, nil
}
