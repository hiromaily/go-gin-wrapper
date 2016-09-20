package bases

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	"github.com/hiromaily/go-gin-wrapper/libs/cors"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	lg "github.com/hiromaily/golibs/log"
	reg "github.com/hiromaily/golibs/regexp"
)

func getURL(scheme, host string, port int) string {
	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

// IsRefererHostValid is check referer for posted page
func IsRefererHostValid(c *gin.Context, pageFrom string) bool {

	srv := conf.GetConf().Server
	webserverURL := getURL(srv.Scheme, srv.Host, srv.Port)

	//TODO:Add feature that switch https to http easily.
	url := fmt.Sprintf("%s/%s", webserverURL, pageFrom)
	lg.Debugf("expected url: %s", url)

	//http://localhost:9999/login
	lg.Debugf("Referer: %s", c.Request.Header.Get("Referer"))
	//lg.Debugf("Referer: %s", c.Request.Referer())

	//default action
	if url != c.Request.Referer() {
		//Invalid access
		lg.Error("Referer is invalid.")

		//token delete
		sess.DelTokenSession(c)

		//set error
		c.AbortWithError(400, errors.New("Referer is invalid."))
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
	//HTTP/1.1
	return c.Request.Proto
}

/*
// SetRquestHeaderForSecurity is to set HTTP request header
func SetRequestHeaderForSecurity(c *gin.Context) {
	c.Request.Header.Set("X-Content-Type-Options", "nosniff")
	c.Request.Header.Set("X-XSS-Protection", "1, mode=block")
	c.Request.Header.Set("X-Frame-Options", "deny")
	c.Request.Header.Set("Content-Security-Policy", "default-src 'none'")
	//c.Request.Header.Set("Strict-Transport-Security", "max-age=15768000")
}
*/

// SetResponseHeaderForSecurity is to set HTTP response header
func SetResponseHeaderForSecurity(c *gin.Context) {
	//http://qiita.com/roothybrid7/items/34578037d883c9a99ca8

	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-XSS-Protection", "1, mode=block")
	c.Writer.Header().Set("X-Frame-Options", "deny")
	c.Writer.Header().Set("Content-Security-Policy", "default-src 'none'")
	//c.Writer.Header().Set("Strict-Transport-Security", "max-age=15768000")

	//CORS
	cors.SetHeader(c)

	//c.Writer.WriteHeader()
	//c.Writer.WriteString()
}

// SetResponseHeader is set HTTP response header (TODO: work in progress)
func SetResponseHeader(c *gin.Context, data []map[string]string) {
	for _, header := range data {
		for key, val := range header {
			c.Request.Header.Set(key, val)
		}
	}
}
