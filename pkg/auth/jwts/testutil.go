package jwts

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
)

// GetJWTToken returns jwt token string from parsed json for test use
func GetJWTToken(client *http.Client, urlString string, postData url.Values) (string, error) {
	ajaxHeader := map[string]string{"X-Requested-With": "XMLHttpRequest"}
	keyHeader := map[string]string{"X-Custom-Header-Gin": "key12345"}
	contentType := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	jwtHeaders := []map[string]string{ajaxHeader, keyHeader, contentType}

	req, _ := http.NewRequest("POST", urlString, bytes.NewBuffer([]byte(postData.Encode())))
	httpheader.SetHTTPHeaders(req, jwtHeaders)

	res, err := client.Do(req)
	defer func() {
		if res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.Errorf("status code is not 200: %d", res.StatusCode)
	}

	return GetJWTResponseToken(res)
}

// GetJWTResponseToken returns jwt token from response
func GetJWTResponseToken(res *http.Response) (string, error) {
	type ResponseJWT struct {
		Code  uint8  `json:"code"`
		Token string `json:"token"`
	}

	body, _, err := parseBody(res)
	if err != nil {
		return "", errors.Wrap(err, "fail to call parseResponse()")
	}

	var jwt ResponseJWT
	if err := json.Unmarshal(body, &jwt); err != nil {
		return "", err
	}
	return jwt.Token, nil
}

// parse response body
func parseBody(res *http.Response) ([]byte, int, error) {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}
	return body, res.StatusCode, nil
}
