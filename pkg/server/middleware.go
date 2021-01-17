package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/auth/jwts"
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/cors"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
	str "github.com/hiromaily/go-gin-wrapper/pkg/strings"
)

// Middlewarer interface
type Middlewarer interface {
	RejectSpecificIP() gin.HandlerFunc
	SetMetaData() gin.HandlerFunc
	UpdateUserSession() gin.HandlerFunc
	GlobalRecover() gin.HandlerFunc
	CheckHTTPRefererAndCSRF() gin.HandlerFunc
	CheckHTTPReferer() gin.HandlerFunc
	CheckCSRF() gin.HandlerFunc
	RejectNonHTTPS() gin.HandlerFunc
	CheckHTTPHeader() gin.HandlerFunc
	CheckJWT() gin.HandlerFunc
	CheckCORS() gin.HandlerFunc
}

// server object
type middleware struct {
	// session xxxx
	logger      *zap.Logger
	jwter       jwts.JWTer
	serverConf  *config.Server
	proxyConf   *config.Proxy
	apiConf     *config.API
	corsConf    *config.CORS
	developConf *config.Develop
}

// NewMiddleware returns Server interface
func NewMiddleware(
	logger *zap.Logger,
	jwter jwts.JWTer,
	serverConf *config.Server,
	proxyConf *config.Proxy,
	apiConf *config.API,
	corsConf *config.CORS,
	developConf *config.Develop,
) Middlewarer {
	return &middleware{
		logger:      logger,
		jwter:       jwter,
		serverConf:  serverConf,
		proxyConf:   proxyConf,
		apiConf:     apiConf,
		corsConf:    corsConf,
		developConf: developConf,
	}
}

// RefererURLs key->request url, value->refer url
var RefererURLs = map[string]string{
	"/login": "login",
	"/user":  "user",
}

//-----------------------------------------------------------------------------
// basic
//-----------------------------------------------------------------------------

// TODO:skip static files like (jpg, gif, png, js, css, woff)

// [WIP] RejectSpecificIP rejects IP addresses in blacklist
// TODO: reject all except reverse　proxy ip address.
func (m *middleware) RejectSpecificIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		if httpheader.IsStaticFile(c) {
			c.Next()
			return
		}

		m.logger.Info("RejectSpecificIp")
		ip := c.ClientIP()
		// proxy
		if m.proxyConf.Mode != 0 {
			if m.proxyConf.Server.Host != ip {
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
func (m *middleware) SetMetaData() gin.HandlerFunc {
	return func(c *gin.Context) {
		if httpheader.IsStaticFile(c) {
			c.Next()
			return
		}
		m.logger.Info("SetMetaData")

		// Context Meta Data
		// http.Header{
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

		// Ajax
		if m.isXHR(c) {
			c.Set("ajax", "1")
		} else {
			c.Set("ajax", "0")
		}

		// Response Data
		if m.isAcceptHeaderJSON(c) {
			c.Set("responseData", "json")
		} else {
			c.Set("responseData", "html")
		}

		// Requested Data
		if m.isContentTypeJSON(c) {
			c.Set("requestData", "json")
		} else {
			c.Set("requestData", "data")
		}

		// User Agent
		c.Set("userAgent", m.getUserAgent(c))

		// Language
		c.Set("language", m.getLanguage(c))

		c.Next()
	}
}

// UpdateUserSession is update user session expire
// TODO:When session has already started, update session expired
func (m *middleware) UpdateUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		if httpheader.IsStaticFile(c) {
			c.Next()
			return
		}
		m.logger.Info("UpdateUserSession")
		if bRet, uid := sess.IsLogin(c); bRet {
			sess.SetUserSession(c, uid)
		}
		c.Next()
	}
}

// GlobalRecover is after request, handle aborted code or 500 error.
// When 404 or 405 error occurred, response already been set in controller/errors/errors.go
func (m *middleware) GlobalRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if httpheader.IsStaticFile(c) {
				c.Next()
				return
			}
			m.logger.Info("[GlobalRecover] defer func()")

			// when crossing request, context data can't be left.
			// c.Set("skipMiddleWare", "1")

			if c.IsAborted() {
				m.logger.Debug("GlobalRecover", zap.Bool("c.IsAborted()", c.IsAborted()))

				// response
				m.setResponse(c, getErrMsg(c), c.Writer.Status())
				return
			}

			if m.developConf.RecoverEnable {
				if rec := recover(); rec != nil {
					m.logger.Debug("GlobalRecover", zap.Any("recover()", rec))
					// TODO:How should response data be decided whether html or json?
					// TODO:Ajax or not doesn't matter to response. HTTP header of Accept may be better.
					// TODO:How precise should I follow specifications of HTTP header.
					m.setResponse(c, str.Itos(rec), 500)
					return
				}
			}
		}(c)

		c.Next()
		// Next is [Main gin Recovery]
	}
}

func (m *middleware) setResponse(c *gin.Context, errMsg string, code int) {
	refURL := "/"
	if c.Request.Header.Get("Referer") != "" {
		refURL = c.Request.Header.Get("Referer")
	}

	if m.isXHR(c) {
		c.JSON(c.Writer.Status(), gin.H{
			"code":  fmt.Sprintf("%d", code),
			"error": errMsg,
		})
	} else {
		c.HTML(c.Writer.Status(), "pages/errors/error.tmpl", gin.H{
			"code":    fmt.Sprintf("%d", code),
			"message": errMsg,
			"url":     refURL,
		})
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
			errMsg = fmt.Sprintf("unexpected error: http status: %d", c.Writer.Status())
		}
	}

	return errMsg
}

//-----------------------------------------------------------------------------
// required by each controller
//-----------------------------------------------------------------------------

// [WIP] CheckHTTPRefererAndCSRF checks referer and CSRF token
func (m *middleware) CheckHTTPRefererAndCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("[CheckHTTPRefererAndCSRF]")
		// Referer
		url := httpheader.GetURL(c)
		// get referer data mapping table using url (map[string])
		if refURL, ok := RefererURLs[url]; ok {
			// Check Referer
			if !httpheader.IsRefererHostValid(c, m.logger, m.serverConf, refURL) {
				c.Next()
				return
			}
		}

		// CSRF
		sess.IsTokenSessionValid(c, m.logger, c.PostForm("gintoken"))
		c.Next()
	}
}

// [WIP] CheckHTTPReferer is to check HTTP Referer.
func (m *middleware) CheckHTTPReferer() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("[heckHttpReferer]")
		url := httpheader.GetURL(c)
		// get referer data mapping table using url (map[string])
		if refURL, ok := RefererURLs[url]; ok {
			// Check Referer
			httpheader.IsRefererHostValid(c, m.logger, m.serverConf, refURL)
		}
		c.Next()
	}
}

// [WIP] CheckCSRF is to check CSRF token
func (m *middleware) CheckCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("[CheckCSRF]")
		sess.IsTokenSessionValid(c, m.logger, c.PostForm("gintoken"))
		c.Next()
	}
}

// [WIP] RejectNonHTTPS is to reject non HTTPS access.
func (m *middleware) RejectNonHTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("[RejectNonHTTPS]")

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
func (m *middleware) CheckHTTPHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("CheckHttpHeader",
			zap.Any("request_header", c.Request.Header),
			zap.Any("request_form", c.Request.Form),
			zap.Any("request_body", c.Request.Body),
			zap.String("client_ip", c.ClientIP()),
			zap.String("request_method", c.Request.Method),
		)
		// IsAjax := c.Request.Header.Get("X-Requested-With")
		// logger.Debugf("[X-Requested-With] %s", IsAjax)

		// IsKey := c.Request.Header.Get("X-Custom-Header-Gin")
		// logger.Debugf("[X-Custom-Header-Gin] %s", IsKey)

		// TODO:when preflight request, X-Requested-With may be not sent
		// TODO:all cors requests don't include X-Requested-With..
		if c.Request.Method != "OPTIONS" && c.Request.Header.Get("X-Custom-Header-Cors") == "" {
			if m.apiConf.Ajax && !m.isXHR(c) {
				// error
				m.logger.Error("Ajax request is required")
				c.AbortWithStatus(400)
				return
			}
		}

		if m.apiConf.Header.Enabled {
			valOfaddedHeader := c.Request.Header.Get(m.apiConf.Header.Header)
			if valOfaddedHeader != m.apiConf.Header.Key {
				// error
				m.logger.Error("header and key are invalid")
				c.AbortWithStatus(400)
				return
			}
		}

		//TODO:if no Content-Type, how should be handled.
		//contentType := c.Request.Header.Get("Content-Type")
		//if contentType == "application/json" {
		//}

		// Context Meta Data
		// SetMetaData(c)

		c.Next()
	}
}

// CheckJWT is to check JWT token code
func (m *middleware) CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("[CheckJWT]")

		IsAuth := c.Request.Header.Get("Authorization")
		if IsAuth == "" {
			c.AbortWithError(400, errors.New("authorization header is missing"))
			return
		}

		var err error
		aAry := strings.Split(IsAuth, " ")
		if len(aAry) != 2 {
			err = errors.New("Authorization header is invalid")
		} else {
			if aAry[0] != "Bearer" {
				err = errors.New("Authorization header is invalid")
			} else {
				token := aAry[1]
				err = m.jwter.ValidateToken(token)
			}
		}
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		c.Next()
	}
}

// CheckCORS is to check CORS
func (m *middleware) CheckCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("CheckCORS")

		if c.Request.Method == "OPTIONS" && c.Request.Header.Get("Origin") != "" {
			cors.CheckHeader(c, m.logger, m.corsConf)
		}

		c.Next()
	}
}

//-----------------------------------------------------------------------------
// functions
//-----------------------------------------------------------------------------

// IsXHR is whether request is Ajax or not
func (m *middleware) isXHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}

// IsAcceptHeaderJSON is whether request require JSON or not
func (m *middleware) isAcceptHeaderJSON(c *gin.Context) bool {
	accept := strings.ToLower(c.Request.Header.Get("Accept"))
	return strings.Contains(accept, "application/json")
}

// IsContentTypeJSON is whether data format of request is JSON or not
func (m *middleware) isContentTypeJSON(c *gin.Context) bool {
	accept := strings.ToLower(c.Request.Header.Get("Content-Type"))
	return strings.Contains(accept, "application/json")
}

// GetUserAgent is get user agent
// TODO: work in progress
func (m *middleware) getUserAgent(c *gin.Context) string {
	// "User-Agent":[]string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36
	// (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36"},
	tmpUserAgent := c.Request.Header.Get("User-Agent")
	return tmpUserAgent
}

// GetLanguage is to get Language of top priority
// TODO: work in progress
func (m *middleware) getLanguage(c *gin.Context) string {
	// "Accept-Language":[]string{"ja,en-US;q=0.8,en;q=0.6,de;q=0.4,nl;q=0.2"},
	tmpLanguage := c.Request.Header.Get("Accept-Language")
	if tmpLanguage != "" {
		return strings.Split(tmpLanguage, ",")[0]
	}
	return ""
}
