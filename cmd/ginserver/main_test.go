package main

import (
	"fmt"
	lg "github.com/hiromaily/golibs/log"
	//u "github.com/hiromaily/golibs/utils"
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

//TODO:change toml settings automatically as program code.
//TODO:test for jwt, first code have to be got.

var r *gin.Engine

// Test Data
var redirectErr error = errors.New("redirect")

var (
	configFiles = []string{"settings.default.toml", "settings.toml", "docker.toml", "heroku.toml", "travis.toml"}
)

var (
	ajaxHeader       = map[string]string{"X-Requested-With": "XMLHttpRequest"}
	keyHeader        = map[string]string{"X-Custom-Header-Gin": "key12345"}
	keyHeaderWrong   = map[string]string{"X-Custom-Header-Gin": "mistake"}
	basicAuthHeaders = map[string]string{"Authorization": "Basic d2ViOnRlc3Q="}
	contentType      = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	refererLogin     = map[string]string{"Referer": "http://hiromaily.com:9999/login"}
	loginHeaders     = []map[string]string{contentType, refererLogin}
	rightHeaders     = []map[string]string{ajaxHeader, keyHeader}
	wrongKeyHeaders  = []map[string]string{ajaxHeader, keyHeaderWrong}
	onlyAjaxHeaders  = []map[string]string{ajaxHeader}
	onlyKeyHeaders   = []map[string]string{keyHeader}
)

var getTests = []struct {
	url      string
	code     int
	method   string
	headers  []map[string]string
	nextPage string
	err      error
}{
	{"/", http.StatusOK, "GET", nil, "", nil},
	{"/index", http.StatusMovedPermanently, "GET", nil, "/", redirectErr},
	{"/login", http.StatusOK, "GET", nil, "", nil},
	{"/logout", http.StatusNotFound, "GET", nil, "", nil},
	{"/news/", http.StatusOK, "GET", nil, "", nil},
	{"/accounts/", http.StatusTemporaryRedirect, "GET", nil, "/login", redirectErr},
	{"/admin/", http.StatusUnauthorized, "GET", nil, "", nil},
	{"/admin/", http.StatusOK, "GET", []map[string]string{basicAuthHeaders}, "", nil},
}

// Test Data for Login
var loginTests = []struct {
	url      string
	code     int
	method   string
	headers  []map[string]string
	nextPage string
	email    string
	pass     string
	tokenFlg bool
	err      error
}{
	//1.access by GET first
	{"/login", http.StatusOK, "GET", nil, "/login", "", "", false, nil},
	// access by POST, but no form data.
	{"/login", http.StatusOK, "POST", loginHeaders, "/login", "", "", true, nil},
	//2.access by GET again
	{"/login", http.StatusOK, "GET", nil, "/login", "", "", false, nil},
	// access by POST, but no email.
	{"/login", http.StatusOK, "POST", loginHeaders, "/login", "", "password", true, nil},
	//3.access by GET again
	{"/login", http.StatusOK, "GET", nil, "/login", "", "", false, nil},
	// access by POST, but no password.
	{"/login", http.StatusOK, "POST", loginHeaders, "/login", "aaaa@test.nl", "", true, nil},
	//4.access by GET again
	{"/login", http.StatusOK, "GET", nil, "/login", "aaaa@test.nl", "", false, nil},
	// access by POST, but invalid email.
	{"/login", http.StatusOK, "POST", loginHeaders, "/login", "abcimail.com", "password", true, nil},
	//5.access by GET again
	{"/login", http.StatusOK, "GET", nil, "/login", "aaaa@test.nl", "", false, nil},
	// access by POST, but shorter password.
	{"/login", http.StatusOK, "POST", loginHeaders, "/login", "aaaa@test.de", "123", true, nil},
	//6.access by GET again
	{"/login", http.StatusOK, "GET", nil, "/login", "", "", false, nil},
	// access by POST, but wrong form data.
	{"/login", http.StatusOK, "POST", loginHeaders, "/login", "aaaa@test.nl", "password", true, nil},
	//7.access by GET again
	{"/login", http.StatusOK, "GET", nil, "/login", "", "", false, nil},
	// access by POST with right data, but no token
	{"/login", http.StatusBadRequest, "POST", loginHeaders, "/login", "aaaa@test.jp", "password", false, nil},
	//8.access by GET again
	{"/login", http.StatusOK, "GET", nil, "/login", "", "", false, nil},
	// access by POST with right data. expect to access next page.
	{"/login", http.StatusFound, "POST", loginHeaders, "/accounts", "aaaa@test.jp", "password", true, redirectErr},
}

// Test Data for ajax API (Not use JWT Auth)
var userId int = 12
var getApiTests = []struct {
	url     string
	code    int
	method  string
	headers []map[string]string
	err     error
}{
	{"/api/users", http.StatusOK, "GET", rightHeaders, nil},
	{"/api/users", http.StatusBadRequest, "GET", wrongKeyHeaders, nil},
	{"/api/users", http.StatusBadRequest, "GET", onlyAjaxHeaders, nil},
	{"/api/users", http.StatusBadRequest, "GET", onlyKeyHeaders, nil},
	{"/api/users", http.StatusBadRequest, "GET", nil, nil},
	{"/api/users", http.StatusBadRequest, "POST", rightHeaders, nil}, //TODO:value is necessary
	{"/api/users", http.StatusNotFound, "PUT", rightHeaders, nil},
	{"/api/users", http.StatusNotFound, "DELETE", rightHeaders, nil},
	{fmt.Sprintf("/api/users/%d", userId), http.StatusOK, "GET", rightHeaders, nil},
	{fmt.Sprintf("/api/users/%d", userId), http.StatusNotFound, "POST", rightHeaders, nil},
	{fmt.Sprintf("/api/users/%d", userId), http.StatusBadRequest, "PUT", rightHeaders, nil}, //TODO:value is necessary
	{fmt.Sprintf("/api/users/%d", userId), http.StatusOK, "DELETE", rightHeaders, nil},
	{fmt.Sprintf("/api/users/%d", userId), http.StatusOK, "GET", rightHeaders, nil}, //TODO:no resource is right
}

//-----------------------------------------------------------------------------
// Test Framework
//-----------------------------------------------------------------------------
// Initialize
func init() {
	//flag.Parse()

	//when changing loglevel
	//lg.InitializeLog(lg.DEBUG_STATUS, lg.LOG_OFF_COUNT, 0, "[GOWEB]", "/var/log/go/test.log")
	//lg.InitializeLog(lg.INFO_STATUS, lg.LOG_OFF_COUNT, 0, "[GOWEB]", "/var/log/go/test.log")
	lg.InitializeLog(lg.WARNING_STATUS, lg.LOG_OFF_COUNT, 0, "[GOWEB]", "/var/log/go/test.log")
}

func setup() {
	//Create test database and test data.
	InitDatabase(1)

	//Server On
	r = SetHTTPServer(1, "../../")
}

func teardown() {
	//TODO:Drop test database and test data.
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

//-----------------------------------------------------------------------------
// functions
//-----------------------------------------------------------------------------
// Create Send Data
func createSendData(email, pass, ginToken string) url.Values {
	data := make(url.Values)

	data.Add("inputEmail", email)
	data.Add("inputPassword", pass)
	if ginToken != "" {
		data.Add("gintoken", ginToken)
	}

	return data
}

// Parse Response
func parseResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return string(contents), res.StatusCode
}

// Set HTTP Header
func setHttpHeaders(req *http.Request, headers []map[string]string) {
	//req.Header.Set("Authorization", "Bearer access-token")
	for _, header := range headers {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
}

// Get Cookie
func getCookies(cookies []string, key string) (val string) {
	for _, cookie := range cookies {
		tmp := strings.Split(cookie, ";")
		tmp = strings.Split(tmp[0], "=")
		if tmp[0] == key {
			val = tmp[1]
			break
		}
	}
	return
}

// Get Cookie
func getCookies2(strURL, key string, jar *cookiejar.Jar) (val string) {
	setCookieURL, _ := url.Parse(strURL)
	cookies := jar.Cookies(setCookieURL) //cookies []*http.Cookie

	fmt.Printf("cookies: %v\n", cookies)
	for _, cookie := range cookies {
		if cookie.Name == key {
			val = cookie.Value
			break
		}
	}
	return
}

// check sent http headers
func checkHttpHeader(req *http.Request) {
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf("[checkHttpHeader] error: %s\n", err)
	} else {
		fmt.Printf("[checkHttpHeader] headers:\n%s\n", b)
	}

	//POST /login HTTP/1.1
	//Host: 127.0.0.1:63513
	//User-Agent: Go-http-client/1.1
	//Content-Length: 0
	//Content-Type: application/x-www-form-urlencoded
	//Cookie: go-web-ginserver=MTQ3MTA1MDQ3MnxOd3dBTkVOQlJGZE1WRTlRVmxoWldGbEVSVTFYVGxKSk5VZFhXalJYVkRWRlNWazJWRnBQVUVWWlJGSklOMUZSUTB0TE0waGFRVkU9fC_7LJ1pOXIOZo8ZXg-R4oO1LFXaSqJtvA3l0f6Qk9DA
	//Referer: http://hiromaily.com:9999/login
	//Accept-Encoding: gzip
}

//
func getToken(res *http.Response) (ret string) {
	//<input type="hidden" name="gintoken" value="{{ .gintoken }}">
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		fmt.Println("[getToken]", err)
		return
	}
	doc.Find("input[name=gintoken]").Each(func(_ int, s *goquery.Selection) {
		if val, ok := s.Attr("value"); ok {
			ret = val
		}
	})

	return
}

//-----------------------------------------------------------------------------
// Test
//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
// Config
//-----------------------------------------------------------------------------
func TestConfigs(t *testing.T) {
	//configFiles
	path := "../../configs/"
	for _, v := range configFiles {
		conf.New(path + v)
	}

	//
	conf.New(*tomlPath)
}

//-----------------------------------------------------------------------------
// Get Request for HTML
//-----------------------------------------------------------------------------
/*
func TestGetRequestOne(t *testing.T) {
	//request
	ts := httptest.NewServer(r)
	defer ts.Close()

	url := "/"
	res, err := http.Get(ts.URL + url)
	if err != nil {
		t.Errorf("[%s] this page can't be access. \n error is %s", url, err)
	} else {
		if _, code := parseResponse(res); code != http.StatusOK {
			t.Errorf("[%s] this page can't be access. \n return code is %d", url, code)
		}
	}
}
*/

// Table driven test
// - request code, redirect and address
func TestGetRequestOnTable(t *testing.T) {
	//request
	ts := httptest.NewServer(r)
	defer ts.Close()

	client := &http.Client{
		Timeout: time.Duration(3) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirect")
		},
	}

	for i, tt := range getTests {
		fmt.Printf("%d [%s] %s\n", i+1, tt.method, ts.URL+tt.url)

		req, _ := http.NewRequest(tt.method, ts.URL+tt.url, nil)
		//Set Http Headers
		if tt.headers != nil {
			setHttpHeaders(req, tt.headers)
		}
		res, err := client.Do(req)
		//res, err := client.Get(ts.URL + tt.url)

		//t.Logf("%#v", err)
		//&url.Error{Op:"Get", URL:"/", Err:(*errors.errorString)(0xc8202101b0)}
		urlError, isUrlErr := err.(*url.Error)
		if isUrlErr && urlError.Err.Error() != tt.err.Error() {
			t.Errorf("[%s] this page can't be access. \n error is %s", tt.url, urlError.Err)
		} else {
			//check expected status code
			if res.StatusCode != tt.code {
				t.Errorf("[%s] status code is not correct. \n return code is %d / expected %d", tt.url, res.StatusCode, tt.code)
			}
		}
		//check next page
		if isUrlErr && urlError.Err.Error() == redirectErr.Error() {
			//t.Log(res.Header["Location"])
			if tt.nextPage != res.Header["Location"][0] {
				t.Errorf("[%s] redirect url is not correct. \n url is %s / expected %s", tt.url, res.Header["Location"][0], tt.nextPage)
			}
		}
		res.Body.Close()
	}
}

//-----------------------------------------------------------------------------
// Login Test
//-----------------------------------------------------------------------------
func TestLogin(t *testing.T) {
	//request
	ts := httptest.NewServer(r)
	defer ts.Close()

	var ginToken string
	var sendData url.Values

	//cookie
	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		Timeout: time.Duration(3) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirect")
		},
		Jar: jar,
	}

	for i, tt := range loginTests {
		fmt.Printf("%d [%s] %s\n", i+1, tt.method, ts.URL+tt.url)

		//Create Post Data
		if tt.method == "POST" {
			//ginToken
			if !tt.tokenFlg {
				ginToken = ""
			}
			sendData = createSendData(tt.email, tt.pass, ginToken)
		} else {
			sendData = nil
		}

		//
		req, _ := http.NewRequest(tt.method, ts.URL+tt.url, bytes.NewBuffer([]byte(sendData.Encode())))

		//Set Http Headers
		if tt.headers != nil {
			setHttpHeaders(req, tt.headers)
		}

		res, err := client.Do(req)

		urlError, isUrlErr := err.(*url.Error)
		if isUrlErr && urlError.Err.Error() != tt.err.Error() {
			t.Errorf("[%s] this page can't be access. \n error is %s", tt.url, urlError.Err)
		} else {
			//check expected status code
			if res.StatusCode != tt.code {
				t.Errorf("[%s] status code is not correct. \n return code is %d / expected %d", tt.url, res.StatusCode, tt.code)
			}
		}
		//check next page
		if isUrlErr && urlError.Err.Error() == redirectErr.Error() {
			//t.Log(res.Header["Location"])
			if tt.nextPage != res.Header["Location"][0] {
				t.Errorf("[%s] redirect url is not correct. \n url is %s / expected %s", tt.url, res.Header["Location"][0], tt.nextPage)
			}
		}

		//get token for next request
		ginToken = getToken(res)

		//cookie
		//val := getCookies(res.Header["Set-Cookie"], "go-web-ginserver")
		//fmt.Printf("go-web-ginserver: \n%s\n", val)

		//check requested http header
		// As you know, cookie is sent without intentional addition
		//checkHttpHeader(req)

		// Close body
		res.Body.Close()
	}
}

//-----------------------------------------------------------------------------
// Get Request for API (Ajax)
//-----------------------------------------------------------------------------
func TestGetAPIRequestOnTable(t *testing.T) {
	//request
	ts := httptest.NewServer(r)
	defer ts.Close()

	client := &http.Client{
		Timeout: time.Duration(3) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirect")
		},
	}

	for i, tt := range getApiTests {
		fmt.Printf("%d [%s] %s\n", i+1, tt.method, ts.URL+tt.url)

		req, _ := http.NewRequest(tt.method, ts.URL+tt.url, nil)
		//Set Http Headers
		if tt.headers != nil {
			setHttpHeaders(req, tt.headers)
		}
		res, err := client.Do(req)

		urlError, isUrlErr := err.(*url.Error)
		if isUrlErr && urlError.Err.Error() != tt.err.Error() {
			t.Errorf("[%s] this page can't be access. \n error is %s", tt.url, urlError.Err)
		} else {
			//check expected status code
			if res.StatusCode != tt.code {
				t.Logf("%#v", tt)
				t.Errorf("[%s] status code is not correct. \n return code is %d / expected %d", tt.url, res.StatusCode, tt.code)
			}
		}

		res.Body.Close()
	}
}
