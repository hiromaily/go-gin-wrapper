package json

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
)

// ResponseUserJSON is Return user json
func ResponseUserJSON(c *gin.Context, co *config.CORSConfig, code int, obj interface{}) {
	// Set Header
	hh.SetResponseHeaderForSecurity(c, co)

	if code == 0 {
		code = http.StatusOK
	}

	//c.JSON(http.StatusOK, gin.H{
	//	"message": "msg",
	//	"name":    "name",
	//})

	c.JSON(code, obj)
}
