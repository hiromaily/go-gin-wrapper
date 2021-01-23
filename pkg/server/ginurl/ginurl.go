package ginurl

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetURLString formats ctx.Request.URL
func GetURLString(ctx *gin.Context) string {
	//&url.URL{
	// Scheme:"",
	// Opaque:"",
	// User:(*url.Userinfo)(nil),
	// Host:"",
	// Path:"/login", RawPath:"", RawQuery:"", Fragment:""}
	ctxURL := ctx.Request.URL
	return fmt.Sprintf("%s%s%s", ctxURL.Scheme, ctxURL.Host, ctxURL.Path)
}
