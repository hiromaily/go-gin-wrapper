package httpheader

import (
	"github.com/gin-gonic/gin"
)

// SetResponseHeader sets HTTP response header
func SetResponseHeader(ctx *gin.Context) {
	ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.Writer.Header().Set("X-XSS-Protection", "1, mode=block")
	ctx.Writer.Header().Set("X-Frame-Options", "deny")
	ctx.Writer.Header().Set("Content-Security-Policy", "default-src 'none'")
	// ctx.Writer.Header().Set("Strict-Transport-Security", "max-age=15768000")
}

// setResponseHeader sets HTTP response header
//func setResponseHeader(ctx *gin.Context, data []map[string]string) {
//	for _, header := range data {
//		for key, val := range header {
//			ctx.Writer.Header().Set(key, val)
//		}
//	}
//}
