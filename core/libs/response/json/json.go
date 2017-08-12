package json

import (
	//"encoding/json"
	hh "github.com/hiromaily/go-gin-wrapper/core/libs/httpheader"
	//gin "gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RtnUserJSON is Return user json
func RtnUserJSON(c *gin.Context, code int, obj interface{}) {
	//Set Header
	hh.SetResponseHeaderForSecurity(c)

	if code == 0 {
		code = http.StatusOK
	}

	//c.JSON(http.StatusOK, gin.H{
	//	"message": "msg",
	//	"name":    "name",
	//})

	c.JSON(code, obj)
}
