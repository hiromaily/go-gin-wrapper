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
	"github.com/hiromaily/go-gin-wrapper/pkg/reverseproxy/types"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/cors"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
	str "github.com/hiromaily/go-gin-wrapper/pkg/strings"
)

// Middlewarer interface
type Middlewarer interface {
	GlobalRecover() gin.HandlerFunc
	FilterIP() gin.HandlerFunc
	SetMetaData() gin.HandlerFunc
	UpdateUserSession() gin.HandlerFunc
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
	rejectIPs   []string
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
	rejectIPs []string,
	serverConf *config.Server,
	proxyConf *config.Proxy,
	apiConf *config.API,
	corsConf *config.CORS,
	developConf *config.Develop,
) Middlewarer {
	return &middleware{
		logger:      logger,
		jwter:       jwter,
		rejectIPs:   rejectIPs,
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

// GlobalRecover is after request, handle aborted code or 500 error.
// When 404 or 405 error occurred, response already been set in controller/errors/errors.go
func (m *middleware) GlobalRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			if !preCheck(c) {
				return
			}

			m.logger.Info("middleware GlobalRecover()")
			if c.IsAborted() {
				m.logger.Debug("GlobalRecover", zap.Bool("c.IsAborted()", c.IsAborted()))
				m.setResponse(c, getErrMsg(c), c.Writer.Status())
				return
			}

			if m.developConf.RecoverEnable {
				if rec := recover(); rec != nil {
					m.logger.Debug("GlobalRecover", zap.Any("recover()", rec))
					m.setResponse(c, str.Itos(rec), http.StatusInternalServerError)
					return
				}
			}
		}(c)

		c.Next()
		// Next is `main gin Recovery`
	}
}

func (m *middleware) setResponse(c *gin.Context, errMsg string, code int) {
	referer := c.Request.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	if m.isAcceptHeaderJSON(c) {
		c.JSON(c.Writer.Status(), gin.H{
			"code":  fmt.Sprintf("%d", code),
			"error": errMsg,
		})
		return
	}
	c.HTML(c.Writer.Status(), "pages/errors/error.tmpl", gin.H{
		"code":    fmt.Sprintf("%d", code),
		"message": errMsg,
		"url":     referer,
	})
}

func getErrMsg(c *gin.Context) string {
	if c.Errors != nil && c.Errors.Last() != nil {
		return c.Errors.Last().Err.Error()
	}

	switch c.Writer.Status() {
	case 400:
		return http.StatusText(http.StatusBadRequest)
	case 401:
		return http.StatusText(http.StatusUnauthorized)
	case 403:
		return http.StatusText(http.StatusForbidden)
	case 404:
		return http.StatusText(http.StatusNotFound)
	case 405:
		return http.StatusText(http.StatusMethodNotAllowed)
	case 406:
		return http.StatusText(http.StatusNotAcceptable)
	case 407:
		return http.StatusText(http.StatusProxyAuthRequired)
	case 408:
		return http.StatusText(http.StatusRequestTimeout)
	case 500:
		return http.StatusText(http.StatusInternalServerError)
	}
	return fmt.Sprintf("unexpected error: http status: %d", c.Writer.Status())
}

// FilterIP rejects IP addresses in blacklist
func (m *middleware) FilterIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !preCheck(c) {
			return
		}

		m.logger.Info("middleware FilterIP")
		ip := c.ClientIP()
		// proxy
		if m.proxyConf.Mode != types.NoProxy && m.proxyConf.Server.Host != ip {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		for _, rejectIP := range m.rejectIPs {
			if ip == rejectIP {
				c.AbortWithStatus(403)
				return
			}
		}
		c.Next()
	}
}

// SetMetaData is to set meta data
func (m *middleware) SetMetaData() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !preCheck(c) {
			return
		}

		m.logger.Info("middleware SetMetaData")

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

// UpdateUserSession updates user session expire
func (m *middleware) UpdateUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !preCheck(c) {
			return
		}

		m.logger.Info("middleware UpdateUserSession")
		if logined, uid := sess.IsLogin(c); logined {
			sess.SetUserSession(c, uid)
		}
		c.Next()
	}
}

//-----------------------------------------------------------------------------
// required by each controller
//-----------------------------------------------------------------------------

// CheckHTTPReferer checks referer
func (m *middleware) CheckHTTPReferer() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("middleware CheckHTTPReferer")
		defer c.Next()

		targetURL, found := RefererURLs[httpheader.GetURL(c)]
		if !found {
			return
		}
		// check referer
		httpheader.IsRefererHostValid(c, m.logger, m.serverConf, targetURL)
	}
}

// CheckHTTPReferer checks CSRF token
func (m *middleware) CheckCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("middleware CheckCSRF")
		sess.IsTokenSessionValid(c, m.logger, c.PostForm("gintoken"))
		c.Next()
	}
}

// RejectNonHTTPS rejects if request is NOT HTTPS
func (m *middleware) RejectNonHTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("middleware RejectNonHTTPS")

		if !strings.Contains(httpheader.GetURL(c), "https://") {
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}

//-----------------------------------------------------------------------------
// web API use
//-----------------------------------------------------------------------------

// CheckHTTPHeader checks HTTP Header
func (m *middleware) CheckHTTPHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("middleware CheckHttpHeader",
			zap.Any("request_header", c.Request.Header),
			zap.Any("request_form", c.Request.Form),
			zap.Any("request_body", c.Request.Body),
			zap.String("client_ip", c.ClientIP()),
			zap.String("request_method", c.Request.Method),
			zap.String("http_header_X-Custom-Header-Gin", c.Request.Header.Get("X-Custom-Header-Gin")),
		)

		if m.apiConf.Header.Enabled {
			// X-Custom-Header-Gin
			if c.Request.Header.Get(m.apiConf.Header.Header) != m.apiConf.Header.Key {
				m.logger.Error("header and key are invalid")
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		// TODO: when preflight request, `X-Requested-With` may be not sent
		if c.Request.Method != "OPTIONS" {
			c.Next()
			return
		}
		// TODO: all cors requests don't include `X-Requested-With`
		if m.apiConf.Ajax && !m.isXHR(c) {
			m.logger.Error("Ajax request is required")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Next()
	}
}

// CheckJWT checks JWT token code
func (m *middleware) CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("middleware CheckJWT")

		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("authorization header is missing"))
			return
		}

		// Bearer token
		authParts := strings.Split(auth, " ")
		if len(authParts) != 2 {
			c.AbortWithError(http.StatusBadRequest, errors.New("Authorization header is invalid"))
			return
		}
		if authParts[0] != "Bearer" {
			c.AbortWithError(http.StatusBadRequest, errors.New("Authorization header is invalid"))
			return
		}
		if err := m.jwter.ValidateToken(authParts[1]); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Next()
	}
}

// CheckCORS checks CORS
func (m *middleware) CheckCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.logger.Info("middleware CheckCORS")

		if c.Request.Method == "OPTIONS" && c.Request.Header.Get("Origin") != "" {
			cors.CheckHeader(c, m.logger, m.corsConf)
		}
		c.Next()
	}
}

//-----------------------------------------------------------------------------
// funcs
//-----------------------------------------------------------------------------

func preCheck(c *gin.Context) bool {
	if httpheader.IsStaticFile(c) {
		c.Next()
		return false
	}
	return true
}

// IsXHR is whether request is Ajax or not
func (m *middleware) isXHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}

// IsAcceptHeaderJSON is whether request accepts JSON or not
func (m *middleware) isAcceptHeaderJSON(c *gin.Context) bool {
	accept := strings.ToLower(c.Request.Header.Get("Accept"))
	return strings.Contains(accept, "application/json")
}

// IsContentTypeJSON is whether Content-Type of request is JSON or not
func (m *middleware) isContentTypeJSON(c *gin.Context) bool {
	accept := strings.ToLower(c.Request.Header.Get("Content-Type"))
	return strings.Contains(accept, "application/json")
}

// GetUserAgent returns user agent
func (m *middleware) getUserAgent(c *gin.Context) string {
	return c.Request.Header.Get("User-Agent")
}

// GetLanguage returns language of highest priority
func (m *middleware) getLanguage(c *gin.Context) string {
	// "Accept-Language":[]string{"ja,en-US;q=0.8,en;q=0.6,de;q=0.4,nl;q=0.2"},
	lang := c.Request.Header.Get("Accept-Language")
	if lang == "" {
		return ""
	}
	return strings.Split(lang, ",")[0]
}
