package cors

import (
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	lg "github.com/hiromaily/golibs/log"
	//"strings"
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
func CheckHeader(c *gin.Context) int {
	lg.Info("[cors.CheckHeader]")
	lg.Debug("%v", )

	// check preflight of XMLHttpRequest Level2 XMLHttpRequest
	// 1.check origin
	origin := c.Request.Header.Get("Origin")
	lg.Debugf("Origin header: %v", origin)

	// 2.check header
	header := c.Request.Header.Get("Access-Control-Request-Headers")
	lg.Debugf("Access-Control-Request-Headers header: %v", header)

	// 3.check method
	method := c.Request.Header.Get("Access-Control-Request-Method")
	lg.Debugf("Access-Control-Request-Method header: %v", method)

	/*
		if c.Request.Method == "OPTIONS" {
			//ヘッダーにAuthorizationが含まれていた場合はpreflight成功
			s := c.Request.Header.Get("Access-Control-Request-Headers")
			if strings.Contains(s, "authorization") == true || strings.Contains(s, "Authorization") == true {
				return 204
				//c.Writer.WriteHeader(204)
			}
			//c.Writer.WriteHeader(400)
			return 400
		}
	*/
	return 0
}

// SetHeader is for CORS
func SetHeader(c *gin.Context) bool {
	lg.Info("[cors.SetHeader]")
	if conf.GetConf().API.CORS.Enabled {
		lg.Debugf("c.Request.RemoteAddr: %v", c.Request.RemoteAddr)
		//[::1]:58434

		//Access-Control-Allow-Origin
		// allow from remote addr
		c.Writer.Header().Set(HeaderOrigin, c.Request.RemoteAddr)

		//Access-Control-Allow-Headers
		c.Writer.Header().Set(HeaderHeaders, "Origin, X-Requested-With, Content-Type, Accept, Authorization")

		//Access-Control-Allow-Methods
		c.Writer.Header().Set(HeaderMethods, "GET, POST, PUT, DELETE, OPTIONS")

		//Access-Control-Allow-Credentials
		c.Writer.Header().Set(HeaderCredentials, "false")

		return true
	}
	return false
}
