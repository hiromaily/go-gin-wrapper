package main

import (
	"fmt"
	lg "github.com/hiromaily/golibs/log"
	//u "github.com/hiromaily/golibs/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

var r *gin.Engine

// Test Data
var redirectErr error = errors.New("redirect")
var getTests = []struct {
	url      string
	code     int
	nextPage string
	err      error
}{
	{"/", http.StatusOK, "", nil},
	{"/index", http.StatusMovedPermanently, "/", redirectErr},
	{"/login", http.StatusOK, "", nil},
	{"/logout", http.StatusNotFound, "", nil},
	{"/news/", http.StatusOK, "", nil},
	{"/accounts/", http.StatusTemporaryRedirect, "/login", redirectErr},
	{"/admin/", http.StatusUnauthorized, "", nil},
}

// Test Data for ajax API
var (
	ajaxHeader       = map[string]string{"X-Requested-With": "XMLHttpRequest"}
	keyHeader        = map[string]string{"X-Custom-Header-Gin": "key12345"}
	keyHeaderWrong   = map[string]string{"X-Custom-Header-Gin": "mistake"}
	basicAuthHeaders = map[string]string{"Authorization": "Basic d2ViOnRlc3Q="}
	rightHeaders     = []map[string]string{ajaxHeader, keyHeader}
	wrongKeyHeaders  = []map[string]string{ajaxHeader, keyHeaderWrong}
	onlyAjaxHeaders  = []map[string]string{ajaxHeader}
	onlyKeyHeaders   = []map[string]string{keyHeader}
)
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

//-----------------------------------------------------------------------------
// Test Framework
//-----------------------------------------------------------------------------
// Initialize
func init() {
	//flag.Parse()
	lg.InitializeLog(lg.INFO_STATUS, lg.LOG_OFF_COUNT, 0, "[GOWEB_TEST]", "/var/log/go/test.log")
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
// Test
//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
// Get Request for HTML
//-----------------------------------------------------------------------------
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

	for _, tt := range getTests {
		res, err := client.Get(ts.URL + tt.url)
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

	for _, tt := range getApiTests {
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
