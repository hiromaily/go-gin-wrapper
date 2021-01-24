package ginbinder

import (
	"github.com/gin-gonic/gin"
)

// Bind calls proper bind func
func Bind(ctx *gin.Context, data interface{}) error {
	if val, isExists := ctx.Get("responseData"); isExists {
		if str, ok := val.(string); ok && str == "json" {
			return ctx.BindJSON(data)
		}
	}
	return ctx.Bind(data)
}
