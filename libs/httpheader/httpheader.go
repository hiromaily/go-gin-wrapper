package bases

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	lg "github.com/hiromaily/golibs/log"
)

//Check Referer for posted page
func IsRefererHostValid(c *gin.Context, pageFrom string) bool {

	server := conf.GetConfInstance().Server
	//TODO:Add feature that switch https to http easily.
	url := fmt.Sprintf("%s/%s", server.Referer, pageFrom)
	lg.Debugf("expected url: %s", url)

	//http://localhost:9999/login
	lg.Debugf("Referer: %s", c.Request.Header.Get("Referer"))
	//lg.Debugf("Referer: %s", c.Request.Referer())

	//default action
	if url != c.Request.Referer() {
		//Invalid access
		lg.Debug("Referer is invalid.")

		//token delete
		sess.DelTokenSession(c)

		//set error
		c.AbortWithError(400, errors.New("Referer is invalid."))
		return false
	}
	return true
}

//Get URL
/*
type URL struct {
	Scheme   string
	Opaque   string    // encoded opaque data
	User     *Userinfo // username and password information
	Host     string    // host or host:port
	Path     string
	RawPath  string // encoded path hint (Go 1.5 and later only; see EscapedPath method)
	RawQuery string // encoded query values, without '?'
	Fragment string // fragment for references, without '#'
}
*/
func GetUrl(c *gin.Context) string {
	//&url.URL{
	// Scheme:"",
	// Opaque:"",
	// User:(*url.Userinfo)(nil),
	// Host:"",
	// Path:"/login", RawPath:"", RawQuery:"", Fragment:""}
	url := c.Request.URL
	return url.Scheme + url.Host + url.Path
}

//Get Proto
func GetProto(c *gin.Context) string {
	//HTTP/1.1
	return c.Request.Proto
}

//Set HTTP Request Header
func SetRquestHeaderForSecurity(c *gin.Context) {
	c.Request.Header.Set("X-Content-Type-Options", "nosniff")
	c.Request.Header.Set("X-XSS-Protection", "1, mode=block")
	c.Request.Header.Set("X-Frame-Options", "deny")
	c.Request.Header.Set("Content-Security-Policy", "default-src 'none'")
	//c.Request.Header.Set("Strict-Transport-Security", "max-age=15768000")
}

//Set HTTP Response Header
func SetResponseHeaderForSecurity(c *gin.Context) {
	//http://qiita.com/roothybrid7/items/34578037d883c9a99ca8

	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-XSS-Protection", "1, mode=block")
	c.Writer.Header().Set("X-Frame-Options", "deny")
	c.Writer.Header().Set("Content-Security-Policy", "default-src 'none'")
	//c.Writer.Header().Set("Strict-Transport-Security", "max-age=15768000")

	//c.Writer.WriteHeader()
	//c.Writer.WriteString()
}

//Set HTTP Response Header [???]
func SetResponseHeader(c *gin.Context, data []map[string]string) {
	for _, header := range data {
		for key, val := range header {
			c.Request.Header.Set(key, val)
		}
	}
}
