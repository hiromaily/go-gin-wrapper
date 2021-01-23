package json

import (
	"github.com/gin-gonic/gin"
)

// ResponseUserJSON is Return user json
func ResponseUserJSON(ctx *gin.Context, code int, obj interface{}) {
	// Set Header TODO: move
	// hh.SetResponseHeader(ctx, logger)

	// CORS TODO: move
	//if corsConf.Enabled && ctx.Request.Method == "GET" {
	//	cors.SetHeader(logger, co)(ctx)
	//}

	//c.JSON(http.StatusOK, gin.H{
	//	"message": "msg",
	//	"name":    "name",
	//})

	ctx.JSON(code, obj)
}
