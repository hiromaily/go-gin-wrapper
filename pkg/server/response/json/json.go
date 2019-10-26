package json

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
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
