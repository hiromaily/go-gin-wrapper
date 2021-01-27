// +build notfixintegration

package main

// TODO: refactoring

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
	"github.com/hiromaily/go-gin-wrapper/pkg/token"
)

var (
	r           *gin.Engine
	errRedirect = errors.New("redirect")
	referer     string

	// header
	// keyHeaderWrong   = map[string]string{"X-Custom-Header-Gin": "mistake"}
	// basicAuthHeaders = map[string]string{"Authorization": "Basic d2ViOnRlc3Q="}
	// jwtAuth          = map[string]string{"Authorization": "Bearer %s"}
	// rightHeaders = []map[string]string{ajaxHeader, keyHeader}
	// wrongKeyHeaders = []map[string]string{ajaxHeader, keyHeaderWrong}
	// onlyAjaxHeaders = []map[string]string{ajaxHeader}
	// onlyKeyHeaders  = []map[string]string{keyHeader}
	// rightHeadersWithJWT = []map[string]string{ajaxHeader, keyHeader, jwtAuth}
)

// Test Data for ajax API (When JWT is off)
//var userID = 12
//
//var getUserAPITests = []struct {
//	url     string
//	code    int
//	method  string
//	headers []map[string]string
//	err     error
//}{
//	{"/api/users", http.StatusOK, "GET", rightHeaders, nil},
//	{"/api/users", http.StatusBadRequest, "GET", wrongKeyHeaders, nil},
//	{"/api/users", http.StatusBadRequest, "GET", onlyAjaxHeaders, nil},
//	{"/api/users", http.StatusBadRequest, "GET", onlyKeyHeaders, nil},
//	{"/api/users", http.StatusBadRequest, "GET", nil, nil},
//	{"/api/users", http.StatusBadRequest, "POST", rightHeaders, nil}, // TODO:value is necessary
//	{"/api/users", http.StatusNotFound, "PUT", rightHeaders, nil},
//	{"/api/users", http.StatusNotFound, "DELETE", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusOK, "GET", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusNotFound, "POST", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusBadRequest, "PUT", rightHeaders, nil}, // TODO:value is necessary
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusOK, "DELETE", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusOK, "GET", rightHeaders, nil}, // TODO:no resource is right
//	// TODO:with post data, put data
//}
//
//// Test Data for ajax API (When JWT is on)
//var getUserAPITests2 = []struct {
//	url     string
//	code    int
//	method  string
//	headers []map[string]string
//	err     error
//}{
//	// no jwts token
//	{"/api/users", http.StatusBadRequest, "GET", rightHeaders, nil},
//	{"/api/users", http.StatusBadRequest, "GET", wrongKeyHeaders, nil},
//	{"/api/users", http.StatusBadRequest, "GET", onlyAjaxHeaders, nil},
//	{"/api/users", http.StatusBadRequest, "GET", onlyKeyHeaders, nil},
//	{"/api/users", http.StatusBadRequest, "GET", nil, nil},
//	{"/api/users", http.StatusBadRequest, "POST", rightHeaders, nil}, // TODO:value is necessary
//	{"/api/users", http.StatusNotFound, "PUT", rightHeaders, nil},
//	{"/api/users", http.StatusNotFound, "DELETE", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusBadRequest, "GET", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusNotFound, "POST", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusBadRequest, "PUT", rightHeaders, nil}, // TODO:value is necessary
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusBadRequest, "DELETE", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusBadRequest, "GET", rightHeaders, nil}, // TODO:no resource is right
//	// TODO:with post data, put data
//	// TODO:with jwts token
//}
//
//// Test Data for ajax API (When JWT is on, plus jwts)
//var getUserAPITests3 = []struct {
//	url     string
//	code    int
//	method  string
//	headers []map[string]string
//	err     error
//}{
//	// with jwts token
//	{"/api/users", http.StatusOK, "GET", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusOK, "GET", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusOK, "DELETE", rightHeaders, nil},
//	{fmt.Sprintf("/api/users/id/%d", userID), http.StatusOK, "GET", rightHeaders, nil}, // TODO:no resource is right
//	// TODO:with post data, put data
//}

func setup() {
	// this code is related to main() in main.go
	flag.Parse()

	// config
	conf, err := config.NewConfig(*tomlPath, *isEncrypted)
	if err != nil {
		panic(err)
	}

	// overwrite config by args
	if *portNum != 0 {
		conf.Server.Port = *portNum
	}

	// Referer for test
	referer = fmt.Sprintf(
		"http://%s:%d/login",
		conf.Server.Host,
		conf.Server.Port)

	// server
	regi := NewRegistry(conf, true) // run as test mode
	server := regi.NewServer()
	if r, err = server.Start(); err != nil {
		panic(err)
	}
}

func teardown() {
	// TODO: drop test database
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func createPostData(email, pass, ginToken string) url.Values {
	data := make(url.Values)

	data.Add("inputEmail", email)
	data.Add("inputPassword", pass)
	if ginToken != "" {
		data.Add("gintoken", ginToken)
	}
	return data
}

func TestGetRequest(t *testing.T) {
	t.SkipNow()

	basicAuthHeaders := map[string]string{"Authorization": "Basic d2ViOnRlc3Q="}

	type args struct {
		url     string
		method  string
		headers []map[string]string
	}
	type want struct {
		statusCode int
		nextPage   string
		err        error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "top page",
			args: args{
				url:     "/",
				method:  "GET",
				headers: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "top page 2",
			args: args{
				url:     "/index",
				method:  "GET",
				headers: nil,
			},
			want: want{
				statusCode: http.StatusMovedPermanently,
				nextPage:   "/",
				err:        errRedirect,
			},
		},
		{
			name: "login page",
			args: args{
				url:     "/login",
				method:  "GET",
				headers: nil,
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "logout page",
			args: args{
				url:     "/logout",
				method:  "GET",
				headers: nil,
			},
			want: want{
				statusCode: http.StatusNotFound,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "accounts page",
			args: args{
				url:     "/accounts/",
				method:  "GET",
				headers: nil,
			},
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				nextPage:   "/login",
				err:        errRedirect,
			},
		},
		{
			name: "admin page without basic auth header",
			args: args{
				url:     "/admin/",
				method:  "GET",
				headers: nil,
			},
			want: want{
				statusCode: http.StatusUnauthorized,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "admin page",
			args: args{
				url:     "/admin/",
				method:  "GET",
				headers: []map[string]string{basicAuthHeaders},
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "",
				err:        nil,
			},
		},
	}

	// request
	ts := httptest.NewServer(r)
	defer ts.Close()

	client := &http.Client{
		Timeout: time.Duration(3) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirect")
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("[%s] %s %s", tt.name, tt.args.method, tt.args.url)

			req, _ := http.NewRequest(tt.args.method, fmt.Sprintf("%s%s", ts.URL, tt.args.url), nil)
			if tt.args.headers != nil {
				httpheader.SetHTTPHeaders(req, tt.args.headers)
			}
			// request
			res, err := client.Do(req)
			defer func() {
				if res.Body != nil {
					res.Body.Close()
				}
			}()

			// handle response
			urlError, isURLErr := err.(*url.Error)
			if isURLErr && urlError.Err.Error() != tt.want.err.Error() {
				t.Errorf("%s request: actual error: %v, want error: %v", tt.args.url, err, tt.want.err)
				return
			}
			if isURLErr && urlError.Err.Error() == errRedirect.Error() {
				if tt.want.nextPage != res.Header["Location"][0] {
					t.Errorf("%s request: actual next page: %v, want: %v", tt.args.url, res.Header["Location"][0], tt.want.nextPage)
				}
			}
			if err != nil {
				return
			}
			if res.StatusCode != tt.want.statusCode {
				t.Errorf("%s request: actual status code: %v, want: %v", tt.args.url, res.StatusCode, tt.want.statusCode)
			}
		})
	}
}

func TestLoginRequest(t *testing.T) {
	t.SkipNow()

	contentType := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	refererLogin := map[string]string{"Referer": referer}
	loginHeaders := []map[string]string{contentType, refererLogin}
	type args struct {
		url     string
		method  string
		headers []map[string]string
		email   string
		pass    string
		isToken bool
	}
	type want struct {
		statusCode int
		nextPage   string
		err        error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		//{
		//	name: "access login page",
		//	args: args{
		//		url: "/login",
		//		method: "GET",
		//		headers: nil,
		//		email: "",
		//		pass: "",
		//		isToken: false,
		//	},
		//	want: want{
		//		statusCode: http.StatusOK,
		//		nextPage: "",
		//		err: nil,
		//	},
		//},
		{
			name: "access login page without email, password",
			args: args{
				url:     "/login",
				method:  "POST",
				headers: loginHeaders,
				email:   "",
				pass:    "",
				isToken: true,
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "access login page without email",
			args: args{
				url:     "/login",
				method:  "POST",
				headers: loginHeaders,
				email:   "",
				pass:    "password",
				isToken: true,
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "access login page without password",
			args: args{
				url:     "/login",
				method:  "POST",
				headers: loginHeaders,
				email:   "foobar@gogin.com",
				pass:    "",
				isToken: true,
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "access login page with invalid email",
			args: args{
				url:     "/login",
				method:  "POST",
				headers: loginHeaders,
				email:   "wrong@gogin.com",
				pass:    "password",
				isToken: true,
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "access login page with invalid password",
			args: args{
				url:     "/login",
				method:  "POST",
				headers: loginHeaders,
				email:   "wrong@gogin.com",
				pass:    "hogehoge",
				isToken: true,
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "access login page without token",
			args: args{
				url:     "/login",
				method:  "POST",
				headers: loginHeaders,
				email:   "foobar@gogin.com",
				pass:    "password",
				isToken: false,
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "",
				err:        nil,
			},
		},
		{
			name: "happy path: access login page",
			args: args{
				url:     "/login",
				method:  "POST",
				headers: loginHeaders,
				email:   "foobar@gogin.com",
				pass:    "password",
				isToken: true,
			},
			want: want{
				statusCode: http.StatusOK,
				nextPage:   "/accounts/",
				err:        errRedirect,
			},
		},
	}

	// request
	ts := httptest.NewServer(r)
	defer ts.Close()

	jar, _ := cookiejar.New(nil) // cookie
	client := &http.Client{
		Timeout: time.Duration(3) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirect")
		},
		Jar: jar,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("[%s] %s %s", tt.name, tt.args.method, tt.args.url)

			var (
				ginToken string
				postData url.Values
				err      error
			)

			// get gin token
			if tt.args.isToken {
				ginToken, err = token.GetToken(client, fmt.Sprintf("%s%s", ts.URL, tt.args.url))
				if err != nil {
					t.Fatalf("fail to call token.GetToken() %v", err)
				}
			}
			if tt.args.method == "POST" {
				postData = createPostData(tt.args.email, tt.args.pass, ginToken)
			}

			req, _ := http.NewRequest(tt.args.method, fmt.Sprintf("%s%s", ts.URL, tt.args.url), bytes.NewBuffer([]byte(postData.Encode())))
			if tt.args.headers != nil {
				httpheader.SetHTTPHeaders(req, tt.args.headers)
			}
			// request
			res, err := client.Do(req)
			defer func() {
				if res.Body != nil {
					res.Body.Close()
				}
			}()

			urlError, isURLErr := err.(*url.Error)
			if isURLErr && urlError.Err.Error() != tt.want.err.Error() {
				t.Errorf("%s request: actual error: %v, want error: %v", tt.args.url, err, tt.want.err)
				return
			}
			if isURLErr && urlError.Err.Error() == errRedirect.Error() {
				if tt.want.nextPage != res.Header["Location"][0] {
					t.Errorf("%s request: actual next page: %v, want: %v", tt.args.url, res.Header["Location"][0], tt.want.nextPage)
				}
			}
			if err != nil {
				return
			}
			if res.StatusCode != tt.want.statusCode {
				t.Errorf("%s request: actual status code: %v, want: %v", tt.args.url, res.StatusCode, tt.want.statusCode)
			}
		})
	}
}

func getJWTToken(res *http.Response) (string, error) {
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

func TestJWTAPIRequest(t *testing.T) {
	t.SkipNow()

	ajaxHeader := map[string]string{"X-Requested-With": "XMLHttpRequest"}
	keyHeader := map[string]string{"X-Custom-Header-Gin": "key12345"}
	contentType := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	apiHeaders := []map[string]string{ajaxHeader, keyHeader}
	jwtHeaders := []map[string]string{ajaxHeader, keyHeader, contentType}

	type args struct {
		url     string
		method  string
		headers []map[string]string
		email   string
		pass    string
		// authMode int  //TODO: how to handle // 0:off, 1:HMAC, 2:RSA
	}
	type want struct {
		statusCode int
		err        error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "wrong method",
			args: args{
				url:     "/api/jwts",
				method:  "GET",
				headers: jwtHeaders,
				email:   "foobar@gogin.com",
				pass:    "password",
			},
			want: want{
				statusCode: http.StatusNotFound,
				err:        nil,
			},
		},
		{
			name: "http header without Authorization",
			args: args{
				url:     "/api/jwts",
				method:  "POST",
				headers: apiHeaders,
				email:   "foobar@gogin.com",
				pass:    "password",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				err:        nil,
			},
		},
		{
			name: "password is wrong",
			args: args{
				url:     "/api/jwts",
				method:  "POST",
				headers: jwtHeaders,
				email:   "foobar@gogin.com",
				pass:    "wrong-password",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				err:        nil,
			},
		},
		{
			name: "happy path",
			args: args{
				url:     "/api/jwts",
				method:  "POST",
				headers: jwtHeaders,
				email:   "foobar@gogin.com",
				pass:    "password",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				err:        nil,
			},
		},
	}

	// request
	ts := httptest.NewServer(r)
	defer ts.Close()

	jar, _ := cookiejar.New(nil) // cookie
	client := &http.Client{
		Timeout: time.Duration(3) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("redirect")
		},
		Jar: jar,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("[%s] %s %s", tt.name, tt.args.method, tt.args.url)

			postData := createPostData(tt.args.email, tt.args.pass, "")

			req, _ := http.NewRequest(tt.args.method, fmt.Sprintf("%s%s", ts.URL, tt.args.url), bytes.NewBuffer([]byte(postData.Encode())))
			if tt.args.headers != nil {
				httpheader.SetHTTPHeaders(req, tt.args.headers)
			}
			// request
			res, err := client.Do(req)
			defer func() {
				if res.Body != nil {
					res.Body.Close()
				}
			}()

			urlError, isURLErr := err.(*url.Error)
			if isURLErr && urlError.Err.Error() != tt.want.err.Error() {
				t.Errorf("%s request: actual error: %v, want error: %v", tt.args.url, err, tt.want.err)
				return
			}
			if err != nil {
				return
			}
			if res.StatusCode != tt.want.statusCode {
				t.Errorf("%s request: actual status code: %v, want: %v", tt.args.url, res.StatusCode, tt.want.statusCode)
			}

			// get jwts for next request
			if res.StatusCode != http.StatusOK {
				return
			}
			if _, err := getJWTToken(res); err != nil {
				t.Errorf("fail to parse response: %v", err)
			}
		})
	}
}

//func TestGetUserAPIRequestOnTable(t *testing.T) {
//	t.SkipNow()
//
//	// request
//	ts := httptest.NewServer(r)
//	defer ts.Close()
//
//	client := &http.Client{
//		Timeout: time.Duration(3) * time.Second,
//		CheckRedirect: func(req *http.Request, via []*http.Request) error {
//			return errors.New("redirect")
//		},
//	}
//
//	getAPITestsData := getUserAPITests
//	if *authMode != 0 {
//		// Auth is on
//		if jwtCode == "" {
//			getAPITestsData = getUserAPITests2
//		} else {
//			getAPITestsData = getUserAPITests3
//			// TODO:set JWT to header
//			jwtAuth["Authorization"] = fmt.Sprintf(jwtAuth["Authorization"], jwtCode)
//		}
//	}
//
//	// for i, tt := range getApiTests {
//	for i, tt := range getAPITestsData {
//		fmt.Printf("%d [%s] %s\n", i+1, tt.method, ts.URL+tt.url)
//
//		req, _ := http.NewRequest(tt.method, ts.URL+tt.url, nil)
//		// Set Http Headers
//		if tt.headers != nil {
//			if jwtCode != "" {
//				tt.headers = append(tt.headers, jwtAuth)
//			}
//			setHTTPHeaders(req, tt.headers)
//		}
//		res, err := client.Do(req)
//
//		urlError, isURLErr := err.(*url.Error)
//		if isURLErr && urlError.Err.Error() != tt.err.Error() {
//			t.Errorf("[%s] this page can't be access. \n error is %s", tt.url, urlError.Err)
//		} else {
//			// check expected status code
//			if res.StatusCode != tt.code {
//				t.Logf("%#v", tt)
//				t.Errorf("[%d][%s] status code is not correct. \n return code is %d / expected %d", i+1, tt.url, res.StatusCode, tt.code)
//			}
//		}
//
//		res.Body.Close()
//	}
//}
