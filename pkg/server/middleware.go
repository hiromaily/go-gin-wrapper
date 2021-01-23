package server

import (
	"fmt"
	"net/http"
	"strings"

	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"

	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/auth/jwts"
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/files"
	"github.com/hiromaily/go-gin-wrapper/pkg/reverseproxy/types"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/cors"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/ginurl"
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
	SetResponseHeader() gin.HandlerFunc
	SetCORSHeader() gin.HandlerFunc
}

// server object
type middleware struct {
	// session xxxx
	logger      *zap.Logger
	jwter       jwts.JWTer
	corser      cors.CORSer
	rejectIPs   []string
	serverConf  *config.Server
	proxyConf   *config.Proxy
	apiConf     *config.API
	developConf *config.Develop
}

// NewMiddleware returns Server interface
func NewMiddleware(
	logger *zap.Logger,
	jwter jwts.JWTer,
	corser cors.CORSer,
	rejectIPs []string,
	serverConf *config.Server,
	proxyConf *config.Proxy,
	apiConf *config.API,
	developConf *config.Develop,
) Middlewarer {
	return &middleware{
		logger:      logger,
		jwter:       jwter,
		corser:      corser,
		rejectIPs:   rejectIPs,
		serverConf:  serverConf,
		proxyConf:   proxyConf,
		apiConf:     apiConf,
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
	return func(ctx *gin.Context) {
		defer func(ctx *gin.Context) {
			if !preCheck(ctx) {
				return
			}

			m.logger.Info("middleware GlobalRecover()")
			if ctx.IsAborted() {
				m.logger.Debug("GlobalRecover", zap.Bool("c.IsAborted()", ctx.IsAborted()))
				m.setResponse(ctx, getErrMsg(ctx), ctx.Writer.Status())
				return
			}

			if m.developConf.RecoverEnable {
				if rec := recover(); rec != nil {
					m.logger.Debug("GlobalRecover", zap.Any("recover()", rec))
					m.setResponse(ctx, str.Itos(rec), http.StatusInternalServerError)
					return
				}
			}
		}(ctx)

		ctx.Next()
		// Next is `main gin Recovery`
	}
}

func (m *middleware) setResponse(ctx *gin.Context, errMsg string, code int) {
	referer := ctx.Request.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	if m.isAcceptHeaderJSON(ctx) {
		ctx.JSON(ctx.Writer.Status(), gin.H{
			"code":  fmt.Sprintf("%d", code),
			"error": errMsg,
		})
		return
	}
	ctx.HTML(ctx.Writer.Status(), "pages/errors/error.tmpl", gin.H{
		"code":    fmt.Sprintf("%d", code),
		"message": errMsg,
		"url":     referer,
	})
}

func getErrMsg(ctx *gin.Context) string {
	if ctx.Errors != nil && ctx.Errors.Last() != nil {
		return ctx.Errors.Last().Err.Error()
	}

	switch ctx.Writer.Status() {
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
	return fmt.Sprintf("unexpected error: http status: %d", ctx.Writer.Status())
}

// FilterIP rejects IP addresses in blacklist
func (m *middleware) FilterIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !preCheck(ctx) {
			return
		}

		m.logger.Info("middleware FilterIP")
		ip := ctx.ClientIP()
		// proxy
		if m.proxyConf.Mode != types.NoProxy && m.proxyConf.Server.Host != ip {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		for _, rejectIP := range m.rejectIPs {
			if ip == rejectIP {
				ctx.AbortWithStatus(403)
				return
			}
		}
		ctx.Next()
	}
}

// SetMetaData is to set meta data
func (m *middleware) SetMetaData() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !preCheck(ctx) {
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
		if m.isXHR(ctx) {
			ctx.Set("ajax", "1")
		} else {
			ctx.Set("ajax", "0")
		}

		// Response Data
		if m.isAcceptHeaderJSON(ctx) {
			ctx.Set("responseData", "json")
		} else {
			ctx.Set("responseData", "html")
		}

		// Requested Data
		if m.isContentTypeJSON(ctx) {
			ctx.Set("requestData", "json")
		} else {
			ctx.Set("requestData", "data")
		}

		// User Agent
		ctx.Set("userAgent", m.getUserAgent(ctx))

		// Language
		ctx.Set("language", m.getLanguage(ctx))

		ctx.Next()
	}
}

// UpdateUserSession updates user session expire
func (m *middleware) UpdateUserSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !preCheck(ctx) {
			return
		}

		m.logger.Info("middleware UpdateUserSession")
		if logined, uid := sess.IsLogin(ctx); logined {
			sess.SetUserSession(ctx, uid)
		}
		ctx.Next()
	}
}

//-----------------------------------------------------------------------------
// required by each controller
//-----------------------------------------------------------------------------

// CheckHTTPReferer checks referer
func (m *middleware) CheckHTTPReferer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m.logger.Info("middleware CheckHTTPReferer")

		// FIXME: looks strange, key requires only path
		targetURL, found := RefererURLs[ginurl.GetURLString(ctx)]
		if !found {
			return
		}

		// check referer
		if err := m.validateReferer(ctx, targetURL); err != nil {
			// invalid access
			m.logger.Error("fail to call validateReferer()", zap.Error(err))

			// delete token
			sess.DelTokenSession(ctx)

			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		ctx.Next()
	}
}

// validateReferer validates referer for target page
func (m *middleware) validateReferer(ctx *gin.Context, pageFrom string) error {
	webserverURL := getServerURL(m.serverConf.Scheme, m.serverConf.Host, m.serverConf.Port)
	referer := fmt.Sprintf("%s/%s", webserverURL, pageFrom)

	m.logger.Debug("validateReferer",
		zap.String("expected_referer", referer),
		zap.String("ctx_referer", ctx.Request.Header.Get("Referer")),
	)

	if referer != ctx.Request.Referer() {
		return errors.New("Referer is invalid")
	}
	return nil
}

// CheckCSRF checks CSRF token
func (m *middleware) CheckCSRF() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m.logger.Info("middleware CheckCSRF")
		sess.IsTokenSessionValid(ctx, m.logger, ctx.PostForm("gintoken"))
		ctx.Next()
	}
}

// RejectNonHTTPS rejects if request is NOT HTTPS
func (m *middleware) RejectNonHTTPS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m.logger.Info("middleware RejectNonHTTPS")

		if !strings.Contains(ctx.Request.URL.Scheme, "https://") {
			ctx.AbortWithStatus(403)
			return
		}
		ctx.Next()
	}
}

//-----------------------------------------------------------------------------
// web API use
//-----------------------------------------------------------------------------

// CheckHTTPHeader checks HTTP Header
func (m *middleware) CheckHTTPHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m.logger.Info("middleware CheckHttpHeader",
			zap.Any("request_header", ctx.Request.Header),
			zap.Any("request_form", ctx.Request.Form),
			zap.Any("request_body", ctx.Request.Body),
			zap.String("client_ip", ctx.ClientIP()),
			zap.String("request_method", ctx.Request.Method),
			zap.String("http_header_X-Custom-Header-Gin", ctx.Request.Header.Get("X-Custom-Header-Gin")),
		)

		if m.apiConf.Header.Enabled {
			// X-Custom-Header-Gin
			if ctx.Request.Header.Get(m.apiConf.Header.Header) != m.apiConf.Header.Key {
				m.logger.Error("header and key are invalid")
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		// TODO: when preflight request, `X-Requested-With` may be not sent
		// TODO: all cors requests should not include `X-Requested-With`
		if ctx.Request.Method != "OPTIONS" {
			ctx.Next()
			return
		}
		if m.apiConf.Ajax && !m.isXHR(ctx) {
			m.logger.Error("Ajax request is required")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ctx.Next()
	}
}

// CheckJWT checks JWT token code
func (m *middleware) CheckJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m.logger.Info("middleware CheckJWT")

		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("authorization header is missing"))
			return
		}

		// Bearer token
		authParts := strings.Split(auth, " ")
		if len(authParts) != 2 {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("Authorization header is invalid"))
			return
		}
		if authParts[0] != "Bearer" {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("Authorization header is invalid"))
			return
		}
		if err := m.jwter.ValidateToken(authParts[1]); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx.Next()
	}
}

// CheckCORS checks CORS
func (m *middleware) CheckCORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m.logger.Info("middleware CheckCORS")

		if ctx.Request.Method != "OPTIONS" || ctx.Request.Header.Get("Origin") == "" {
			ctx.Next()
			return
		}
		if err := m.corser.ValidateHeader(ctx); err != nil {
			m.logger.Error("CheckCORS", zap.Error(err))
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
}

// SetResponseHeader sets response header
func (m *middleware) SetResponseHeader() gin.HandlerFunc {
	return hh.SetResponseHeader
}

// SetCORSHeader sets CORS header
func (m *middleware) SetCORSHeader() gin.HandlerFunc {
	return m.corser.SetResponseHeader
}

//-----------------------------------------------------------------------------
// funcs
//-----------------------------------------------------------------------------

func getServerURL(scheme, host string, port int) string {
	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

func preCheck(ctx *gin.Context) bool {
	if files.IsStaticFile(ginurl.GetURLString(ctx)) {
		ctx.Next()
		return false
	}
	return true
}

// IsXHR is whether request is Ajax or not
func (m *middleware) isXHR(ctx *gin.Context) bool {
	return strings.ToLower(ctx.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}

// IsAcceptHeaderJSON is whether request accepts JSON or not
func (m *middleware) isAcceptHeaderJSON(ctx *gin.Context) bool {
	accept := strings.ToLower(ctx.Request.Header.Get("Accept"))
	return strings.Contains(accept, "application/json")
}

// IsContentTypeJSON is whether Content-Type of request is JSON or not
func (m *middleware) isContentTypeJSON(ctx *gin.Context) bool {
	accept := strings.ToLower(ctx.Request.Header.Get("Content-Type"))
	return strings.Contains(accept, "application/json")
}

// GetUserAgent returns user agent
func (m *middleware) getUserAgent(ctx *gin.Context) string {
	return ctx.Request.Header.Get("User-Agent")
}

// GetLanguage returns language of highest priority
func (m *middleware) getLanguage(ctx *gin.Context) string {
	// "Accept-Language":[]string{"ja,en-US;q=0.8,en;q=0.6,de;q=0.4,nl;q=0.2"},
	lang := ctx.Request.Header.Get("Accept-Language")
	if lang == "" {
		return ""
	}
	return strings.Split(lang, ",")[0]
}
