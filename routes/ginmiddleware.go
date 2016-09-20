package routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	"github.com/hiromaily/go-gin-wrapper/libs/cors"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	hh "github.com/hiromaily/go-gin-wrapper/libs/httpheader"
	"github.com/hiromaily/golibs/auth/jwt"
	lg "github.com/hiromaily/golibs/log"
	u "github.com/hiromaily/golibs/utils"
	"net/http"
	"strings"
)

//TODO:skip static files like (jpg, gif, png, js, css, woff)

//-----------------------------------------------------------------------------
// Common
//-----------------------------------------------------------------------------

// RejectSpecificIP is to check registered IP address to reject
// TODO: working in progress yet.
// TODO: it reject all without reverse　proxy ip address.
func RejectSpecificIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		if hh.IsStaticFile(c) {
			c.Next()
			return
		}
		lg.Info("[RejectSpecificIp]")

		ip := c.ClientIP()
		//proxy
		if conf.GetConf().Proxy.Mode != 0 {
			if conf.GetConf().Proxy.Server.Host != ip {
				c.AbortWithStatus(403)
				return
			}
		}
		//if ip != "127.0.0.1" {
		//	c.AbortWithStatus(403)
		//}
		c.Next()
	}
}

// SetMetaData is to set meta data
func SetMetaData() gin.HandlerFunc {
	return func(c *gin.Context) {
		if hh.IsStaticFile(c) {
			c.Next()
			return
		}
		lg.Info("[SetMetaData]")

		//Context Meta Data
		//http.Header{
		// "Referer":[]string{"http://localhost:9999/"},
		// "Accept-Language":[]string{"ja,en-US;q=0.8,en;q=0.6,de;q=0.4,nl;q=0.2"},
		// "X-Frame-Options":[]string{"deny"},
		// "Content-Security-Policy":[]string{"default-src 'none'"},
		// "X-Xss-Protection":[]string{"1, mode=block"},
		// "Connection":[]string{"keep-alive"},
		// "Accept":[]string{"application/json, text/javascript, */*; q=0.01"},
		// "User-Agent":[]string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36"},
		// "X-Content-Type-Options":[]string{"nosniff"},
		// "X-Requested-With":[]string{"XMLHttpRequest"},
		// "X-Custom-Header-Gin":[]string{"key"},
		// "Content-Type":[]string{"application/x-www-form-urlencoded"},
		// "Accept-Encoding":[]string{"gzip, deflate, sdch"}}

		//Ajax
		if IsXHR(c) {
			c.Set("ajax", "1")
		} else {
			c.Set("ajax", "0")
		}

		//Response Data
		if IsAcceptHeaderJSON(c) {
			c.Set("responseData", "json")
		} else {
			c.Set("responseData", "html")
		}

		//Requested Data
		if IsContentTypeJSON(c) {
			c.Set("requestData", "json")
		} else {
			c.Set("requestData", "data")
		}

		//User Agent
		c.Set("userAgent", GetUserAgent(c))

		//Language
		c.Set("language", GetLanguage(c))

		c.Next()
	}
}

// UpdateUserSession is update user session expire
//TODO:When session has already started, update session expired
func UpdateUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		if hh.IsStaticFile(c) {
			c.Next()
			return
		}
		lg.Info("[UpdateUserSession]")
		if bRet, uid := sess.IsLogin(c); bRet {
			sess.SetUserSession(c, uid)
		}
		c.Next()
	}
}

// GlobalRecover is after request, handle aborted code or 500 error.
// When 404 or 405 error occurred, response already been set in controllers/errors/errors.go
func GlobalRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if hh.IsStaticFile(c) {
				c.Next()
				return
			}
			lg.Info("[GlobalRecover] defer func()")

			//when crossing request, context data can't be left.
			//c.Set("skipMiddleWare", "1")

			if c.IsAborted() {
				lg.Debug("[GlobalRecover] c.IsAborted() is true")

				// response
				setResponse(c, getErrMsg(c), c.Writer.Status())
				return
			}

			if conf.GetConf().Develop.RecoverEnable {
				if rec := recover(); rec != nil {
					lg.Debugf("[GlobalRecover] recover() is not nil:\n %v", rec)
					//TODO:How should response data be decided whether html or json?
					//TODO:Ajax or not doesn't matter to response. HTTP header of Accept may be better.
					//TODO:How precise should I follow specifications of HTTP header.

					setResponse(c, u.Itos(rec), 500)
					return
				}
			}
		}(c)

		c.Next()
		//Next is [Main gin Recovery]
	}
}

func getErrMsg(c *gin.Context) string {
	var errMsg string

	if c.Errors != nil {
		if c.Errors.Last() != nil {
			errMsg = c.Errors.Last().Err.Error()
		}
	}

	if errMsg == "" {
		switch c.Writer.Status() {
		case 400:
			errMsg = http.StatusText(http.StatusBadRequest)
		case 401:
			errMsg = http.StatusText(http.StatusUnauthorized)
		case 403:
			errMsg = http.StatusText(http.StatusForbidden)
		case 404:
			errMsg = http.StatusText(http.StatusNotFound)
		case 405:
			errMsg = http.StatusText(http.StatusMethodNotAllowed)
		case 406:
			errMsg = http.StatusText(http.StatusNotAcceptable)
		case 407:
			errMsg = http.StatusText(http.StatusProxyAuthRequired)
		case 408:
			errMsg = http.StatusText(http.StatusRequestTimeout)
		case 500:
			errMsg = http.StatusText(http.StatusInternalServerError)
		default:
			errMsg = "something error is happened."
		}
	}

	return errMsg
}

func setResponse(c *gin.Context, errMsg string, code int) {

	refURL := "/"
	if c.Request.Header.Get("Referer") != "" {
		refURL = c.Request.Header.Get("Referer")
	}

	if IsXHR(c) {
		c.JSON(c.Writer.Status(), gin.H{
			"code": fmt.Sprintf("%d", code),
			//"error": c.Errors.Last().Err.Error(), //it caused error because of nil of objct
			"error": errMsg,
		})
	} else {
		c.HTML(c.Writer.Status(), "pages/errors/error.tmpl", gin.H{
			"code":    fmt.Sprintf("%d", code),
			"message": errMsg,
			"url":     refURL,
		})
	}
	/*
		if IsXHR(c) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  "500",
				"error": rec,
			})
		} else {
			c.HTML(http.StatusInternalServerError, "pages/errors/error.tmpl", gin.H{
				"code":    "500",
				"message": rec,
				"url":     refURL,
			})
		}
	*/

}

//-----------------------------------------------------------------------------
// Respective
//-----------------------------------------------------------------------------

// CheckHTTPRefererAndCSRF is to check referer and CSRF token
// TODO: it's not finished yet.
func CheckHTTPRefererAndCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[CheckHTTPRefererAndCSRF]")
		//Referer
		url := hh.GetURL(c)
		//get referer data mapping table using url (map[string])
		if refURL, ok := RefererURLs[url]; ok {
			//Check Referer
			if !hh.IsRefererHostValid(c, refURL) {
				c.Next()
				return
			}
		}

		//CSRF
		sess.IsTokenSessionValid(c, c.PostForm("gintoken"))
		c.Next()
	}
}

// CheckHTTPReferer is to check HTTP Referer.
// TODO: it's not checked yet if it work well or not.
func CheckHTTPReferer() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[heckHttpReferer]")
		url := hh.GetURL(c)
		//get referer data mapping table using url (map[string])
		if refURL, ok := RefererURLs[url]; ok {
			//Check Referer
			hh.IsRefererHostValid(c, refURL)
		}
		c.Next()
	}
}

// CheckCSRF is to check CSRF token
// TODO: it's not finished yet.
func CheckCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[CheckCSRF]")
		sess.IsTokenSessionValid(c, c.PostForm("gintoken"))
		c.Next()
	}
}

// RejectNonHTTPS is to reject non HTTPS access.
// TODO: it's not fixed yet.
func RejectNonHTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[RejectNonHTTPS]")

		//TODO:check protocol of url
		//if strings.Index(c.url, "https://") == -1 {
		//	c.AbortWithStatus(403)
		//}
		c.Next()
	}
}

//-----------------------------------------------------------------------------
// For Web API
//-----------------------------------------------------------------------------

// CheckHTTPHeader is to check HTTP Header for Ajax request. (For REST)
func CheckHTTPHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[CheckHttpHeader]")

		lg.Debugf("[Request Header]\n%#v\n", c.Request.Header)
		lg.Debugf("[Request Form]\n%#v\n", c.Request.Form)
		lg.Debugf("[Request Body]\n%#v\n", c.Request.Body)
		lg.Debugf("c.ClientIP() %s", c.ClientIP())

		//IsAjax := c.Request.Header.Get("X-Requested-With")
		//lg.Debugf("[X-Requested-With] %s", IsAjax)

		//IsKey := c.Request.Header.Get("X-Custom-Header-Gin")
		//lg.Debugf("[X-Custom-Header-Gin] %s", IsKey)

		apiConf := conf.GetConf().API
		if apiConf.Ajax && !IsXHR(c) {
			//error
			c.AbortWithStatus(400)
			return
		}

		if apiConf.Header.Enabled {
			valOfaddedHeader := c.Request.Header.Get(apiConf.Header.Header)
			if valOfaddedHeader != apiConf.Header.Key {
				//error
				c.AbortWithStatus(400)
				return
			}
		}

		//TODO:if no Content-Type, how should be handled.
		//contentType := c.Request.Header.Get("Content-Type")
		//if contentType == "application/json" {
		//}

		//Context Meta Data
		//SetMetaData(c)

		c.Next()
	}
}

// CheckJWT is to check JWT token code
func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[CheckJWT]")

		var err error

		IsAuth := c.Request.Header.Get("Authorization")
		if IsAuth != "" {
			aAry := strings.Split(IsAuth, " ")
			if len(aAry) != 2 {
				err = errors.New("Authorization header is invalid")
			} else {
				if aAry[0] != "Bearer" {
					err = errors.New("Authorization header is invalid")
				} else {
					token := aAry[1]
					err = jwt.JudgeJWT(token)
				}
			}
		} else {
			err = errors.New("Authorization header was missed.")
		}

		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		c.Next()
	}
}

// CheckCORS is to check CORS
func CheckCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[CheckCORS]")

		cors.CheckHeader(c)

		c.Next()
	}
}

//-----------------------------------------------------------------------------
// functions
//-----------------------------------------------------------------------------

// IsXHR is whether request is Ajax or not
func IsXHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}

// IsAcceptHeaderJSON is whether request require JSON or not
func IsAcceptHeaderJSON(c *gin.Context) bool {
	accept := strings.ToLower(c.Request.Header.Get("Accept"))
	return strings.Index(accept, "application/json") != -1
}

// IsContentTypeJSON is whether data format of request is JSON or not
func IsContentTypeJSON(c *gin.Context) bool {
	accept := strings.ToLower(c.Request.Header.Get("Content-Type"))
	return strings.Index(accept, "application/json") != -1
}

// GetUserAgent is get user agent
// TODO: work in progress
func GetUserAgent(c *gin.Context) string {
	// "User-Agent":[]string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36
	// (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36"},
	tmpUserAgent := c.Request.Header.Get("User-Agent")
	return tmpUserAgent
}

// GetLanguage is to get Language of top priority
// TODO: work in progress
func GetLanguage(c *gin.Context) string {
	// "Accept-Language":[]string{"ja,en-US;q=0.8,en;q=0.6,de;q=0.4,nl;q=0.2"},
	tmpLanguage := c.Request.Header.Get("Accept-Language")
	if tmpLanguage != "" {
		return strings.Split(tmpLanguage, ",")[0]
	}
	return ""
}
