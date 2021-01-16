package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// Adminer interface
type Adminer interface {
	AdminIndexAction(c *gin.Context)
}

// AdminIndexAction [GET]
func (ctl *controller) AdminIndexAction(c *gin.Context) {
	key := c.MustGet(gin.AuthUserKey).(string)
	ctl.logger.Debug("AdminIndexAction", zap.String("gin.AuthUserKey", key))

	// View
	res := gin.H{
		"title":    "Admin Page",
		"navi_key": "/admin/",
	}
	c.HTML(http.StatusOK, "pages/admins/gallery.tmpl", html.Response(res, ctl.apiHeader))
}
