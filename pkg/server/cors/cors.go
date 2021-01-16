package cors

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	str "github.com/hiromaily/go-gin-wrapper/pkg/strings"
)

const (
	// HeaderOrigin is for allowed specific domains
	HeaderOrigin = "Access-Control-Allow-Origin"
	// HeaderHeaders is for allowed headers from browser
	HeaderHeaders = "Access-Control-Allow-Headers"
	// HeaderMethods is for allowed method from browser
	HeaderMethods = "Access-Control-Allow-Methods"
	// HeaderCredentials is whether credentials is used of not
	HeaderCredentials = "Access-Control-Allow-Credentials"
)

// CheckHeader is for CORS before handling request
//  check preflight of XMLHttpRequest Level2 XMLHttpRequest
func CheckHeader(c *gin.Context, logger *zap.Logger, corsConf *config.CORS) {
	logger.Info("CheckHeader")

	// TODO: should it be checked when `Options` method??
	// 1.check origin
	origin := c.Request.Header.Get("Origin")
	logger.Debug("CheckHeader", zap.String("origin header", origin))
	// Origin header: http://127.0.0.1:8000
	if str.SearchIndex(origin, corsConf.Origins) == -1 {
		logger.Error("CheckHeader", zap.Error(errors.New("origin header is invalid")))
		c.AbortWithStatus(400)
		return
	}

	// 2.check header
	// added header intentionally set on Access-Control-Request-Headers
	// this value may change into lower case when extracting
	header := c.Request.Header.Get("Access-Control-Request-Headers")
	logger.Debug("CheckHeader", zap.String("Access-Control-Request-Headers header", header))
	// Access-Control-Request-Headers header: x-custom-header-cors, x-custom-header-gin

	// TODO: is it OK to check
	for _, h := range strings.Split(header, ",") {
		if str.SearchIndexLower(strings.TrimSpace(h), corsConf.Headers) == -1 {
			logger.Error("CheckHeader", zap.Error(errors.New("Access-Control-Request-Headers header is invalid")))
			c.AbortWithStatus(400)
			return
		}
	}

	// 3.check method
	// this value is suppoused request to send after option method request
	method := c.Request.Header.Get("Access-Control-Request-Method")
	logger.Debug("CheckHeader", zap.String("Access-Control-Request-Method header", method))
	// Access-Control-Request-Method header: GET
	if str.SearchIndex(method, corsConf.Methods) == -1 {
		logger.Error("CheckHeader", zap.Error(errors.New("Access-Control-Request-Method header is invalid")))
		c.AbortWithStatus(400)
		return
	}
	// if strings.Contains(s, "authorization") == true || strings.Contains(s, "Authorization") == true {
}

// SetHeader is for CORS
func SetHeader(logger *zap.Logger, corsConf *config.CORS) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("SetHeader")
		if corsConf.Enabled {
			// Access-Control-Allow-Origin
			// allow from remote addr
			logger.Debug("", zap.String("c.Request.RemoteAddr", c.Request.RemoteAddr))
			// c.Writer.Header().Set(HeaderOrigin, c.Request.RemoteAddr)
			c.Writer.Header().Set(HeaderOrigin, strings.Join(corsConf.Origins, ", "))

			// Access-Control-Allow-Headers
			// c.Writer.Header().Set(HeaderHeaders, "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Writer.Header().Set(HeaderHeaders, strings.Join(corsConf.Headers, ", "))

			// Access-Control-Allow-Methods
			// c.Writer.Header().Set(HeaderMethods, "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set(HeaderMethods, strings.Join(corsConf.Methods, ", "))

			// Access-Control-Allow-Credentials
			// c.Writer.Header().Set(HeaderCredentials, "false")
			c.Writer.Header().Set(HeaderCredentials, fmt.Sprintf("%t", corsConf.Credentials))

			return
		}
	}
}
