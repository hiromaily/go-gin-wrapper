package cors

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	lg "github.com/hiromaily/golibs/log"
	u "github.com/hiromaily/golibs/utils"
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
func CheckHeader(c *gin.Context, co *configs.CORSConfig) {
	lg.Info("[cors.CheckHeader]")

	//TODO:Optionsメソッド時にのみチェック？？
	// 1.check origin
	origin := c.Request.Header.Get("Origin")
	lg.Debugf("Origin header: %v", origin)
	//Origin header: http://127.0.0.1:8000
	if u.SearchString(co.Origins, origin) == -1 {
		lg.Error("Origin header is invalid")
		c.AbortWithStatus(400)
		return
	}

	// 2.check header
	// added header intentionally set on Access-Control-Request-Headers
	// this value may change into lower case when extracting
	header := c.Request.Header.Get("Access-Control-Request-Headers")
	lg.Debugf("Access-Control-Request-Headers header: %v", header)
	//Access-Control-Request-Headers header: x-custom-header-cors, x-custom-header-gin

	//TODO:Is it OK to check
	for _, h := range strings.Split(header, ",") {
		if u.SearchStringLower(co.Headers, strings.TrimSpace(h)) == -1 {
			lg.Error("Access-Control-Request-Headers header is invalid")
			c.AbortWithStatus(400)
			return
		}
	}

	// 3.check method
	// this value is suppoused request to send after option method request
	method := c.Request.Header.Get("Access-Control-Request-Method")
	lg.Debugf("Access-Control-Request-Method header: %v", method)
	//Access-Control-Request-Method header: GET
	if u.SearchString(co.Methods, method) == -1 {
		lg.Error("Access-Control-Request-Method header is invalid")
		c.AbortWithStatus(400)
		return
	}
	//if strings.Contains(s, "authorization") == true || strings.Contains(s, "Authorization") == true {
}

// SetHeader is for CORS
func SetHeader(co *configs.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[cors.SetHeader]")
		if co.Enabled {
			//Access-Control-Allow-Origin
			// allow from remote addr
			lg.Debugf("c.Request.RemoteAddr: %v", c.Request.RemoteAddr)
			//c.Writer.Header().Set(HeaderOrigin, c.Request.RemoteAddr)
			c.Writer.Header().Set(HeaderOrigin, strings.Join(co.Origins, ", "))

			//Access-Control-Allow-Headers
			//c.Writer.Header().Set(HeaderHeaders, "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Writer.Header().Set(HeaderHeaders, strings.Join(co.Headers, ", "))

			//Access-Control-Allow-Methods
			//c.Writer.Header().Set(HeaderMethods, "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set(HeaderMethods, strings.Join(co.Methods, ", "))

			//Access-Control-Allow-Credentials
			//c.Writer.Header().Set(HeaderCredentials, "false")
			c.Writer.Header().Set(HeaderCredentials, fmt.Sprintf("%t", co.Credentials))

			return
		}
	}
}
