package html

import (
	"github.com/gin-gonic/gin"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
)

// Response adds common parameter for html response
func Response(obj gin.H, api *config.Header) gin.H {
	// type H map[string]interface{}
	obj["header"] = api.Header
	obj["key"] = api.Key

	return obj
}
