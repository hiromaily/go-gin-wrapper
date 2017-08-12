package html

import (
	conf "github.com/hiromaily/go-gin-wrapper/core/configs"
	//gin "gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-gonic/gin"
)

// Response is to add common parameter for html response
func Response(obj gin.H) gin.H {
	//type H map[string]interface{}

	api := conf.GetConf().API.Header

	obj["header"] = api.Header
	obj["key"] = api.Key

	return obj
}
