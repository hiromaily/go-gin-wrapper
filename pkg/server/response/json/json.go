package json

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
)

// ResponseUserJSON is Return user json
func ResponseUserJSON(c *gin.Context, logger *zap.Logger, corsConf *config.CORS, code int, obj interface{}) {
	// Set Header
	hh.SetResponseHeaderForSecurity(c, logger, corsConf)

	//c.JSON(http.StatusOK, gin.H{
	//	"message": "msg",
	//	"name":    "name",
	//})

	c.JSON(code, obj)
}
