package httpheader

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/cors"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	reg "github.com/hiromaily/golibs/regexp"
)

func getURL(scheme, host string, port int) string {
	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

// IsRefererHostValid is check referer for posted page
func IsRefererHostValid(c *gin.Context, logger *zap.Logger, srvConf *config.ServerConfig, pageFrom string) bool {
	webserverURL := getURL(srvConf.Scheme, srvConf.Host, srvConf.Port)

	// TODO:Add feature that switch https to http easily.
	url := fmt.Sprintf("%s/%s", webserverURL, pageFrom)
	logger.Debug("IsRefererHostValid", zap.String("expected url", url))

	// http://localhost:9999/login
	logger.Debug("IsRefererHostValid", zap.String("Referer", c.Request.Header.Get("Referer")))
	// lg.Debugf("Referer: %s", c.Request.Referer())

	// default action
	if url != c.Request.Referer() {
		// Invalid access
		logger.Error("IsRefererHostValid", zap.Error(errors.New("Referer is invalid")))

		// token delete
		sess.DelTokenSession(c)

		// set error
		c.AbortWithError(400, errors.New("referer is invalid"))
		return false
	}
	return true
}

// GetURL is to get request URL
func GetURL(c *gin.Context) string {
	//&url.URL{
	// Scheme:"",
	// Opaque:"",
	// User:(*url.Userinfo)(nil),
	// Host:"",
	// Path:"/login", RawPath:"", RawQuery:"", Fragment:""}
	url := c.Request.URL
	return url.Scheme + url.Host + url.Path
}

// IsStaticFile is whether request url is for static file or golang page
func IsStaticFile(c *gin.Context) bool {
	url := GetURL(c)
	return reg.IsStaticFile(url)
}

// GetProto is to get protocol
func GetProto(c *gin.Context) string {
	// HTTP/1.1
	return c.Request.Proto
}

// SetResponseHeaderForSecurity is to set HTTP response header
// TODO:it may be better to set config
func SetResponseHeaderForSecurity(c *gin.Context, logger *zap.Logger, co *config.CORSConfig) {
	logger.Info("SetResponseHeaderForSecurity")
	// http://qiita.com/roothybrid7/items/34578037d883c9a99ca8

	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-XSS-Protection", "1, mode=block")
	c.Writer.Header().Set("X-Frame-Options", "deny")
	c.Writer.Header().Set("Content-Security-Policy", "default-src 'none'")
	// c.Writer.Header().Set("Strict-Transport-Security", "max-age=15768000")

	// CORS
	if co.Enabled && c.Request.Method == "GET" {
		cors.SetHeader(logger, co)(c)
	}
	// c.Writer.WriteHeader()
	// c.Writer.WriteString()
}

// SetResponseHeader is set HTTP response header
// TODO: work in progress
func SetResponseHeader(c *gin.Context, data []map[string]string) {
	for _, header := range data {
		for key, val := range header {
			c.Request.Header.Set(key, val)
		}
	}
}
