package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// Adminer interface
type Adminer interface {
	AdminIndexAction(ctx *gin.Context)
}

// AdminIndexAction [GET]
func (ctl *controller) AdminIndexAction(ctx *gin.Context) {
	key := ctx.MustGet(gin.AuthUserKey).(string)
	ctl.logger.Debug("controller AdminIndexAction", zap.String("gin.AuthUserKey", key))

	// view
	res := gin.H{
		"title":    "Admin Page",
		"navi_key": "/admin/",
	}
	ctx.HTML(http.StatusOK, "pages/admins/gallery.tmpl", html.Response(res, ctl.apiHeader))
}
