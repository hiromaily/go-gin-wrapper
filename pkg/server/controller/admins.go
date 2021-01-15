package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// get user, it was set by the BasicAuth middleware

// AdminIndexAction [GET]
func (ctl *Controller) AdminIndexAction(c *gin.Context) {
	key := c.MustGet(gin.AuthUserKey).(string)
	ctl.logger.Debug("AdminIndexAction", zap.String("gin.AuthUserKey", key))

	// View
	res := gin.H{
		"title":    "Admin Page",
		"navi_key": "/admin/",
	}
	c.HTML(http.StatusOK, "pages/admins/gallery.tmpl", html.Response(res, ctl.apiHeader))
}
