package json

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
)

// ResponseUserJSON is Return user json
func ResponseUserJSON(c *gin.Context, logger *zap.Logger, co *config.CORS, code int, obj interface{}) {
	// Set Header
	hh.SetResponseHeaderForSecurity(c, logger, co)

	if code == 0 {
		code = http.StatusOK
	}

	//c.JSON(http.StatusOK, gin.H{
	//	"message": "msg",
	//	"name":    "name",
	//})

	c.JSON(code, obj)
}
