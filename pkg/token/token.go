package token

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Generator interface
type Generator interface {
	Generate() string
}

type generator struct {
	salt string
}

// NewGenerator returns Generator
func NewGenerator(salt string) Generator {
	return &generator{
		salt: salt,
	}
}

func (g *generator) Generate() string {
	md5Hash := md5.New()

	io.WriteString(md5Hash, strconv.FormatInt(time.Now().UnixNano(), 10))
	io.WriteString(md5Hash, g.salt)
	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

// GetToken returns token string from parsed html
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
