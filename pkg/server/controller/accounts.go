package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// Accounter interface
type Accounter interface {
	AccountIndexAction(ctx *gin.Context)
}

// AccountIndexAction [GET]
func (ctl *controller) AccountIndexAction(ctx *gin.Context) {
	ctl.logger.Info("controller AccountIndexAction")

	// validate access
	if logined, _ := sess.IsLogin(ctx); !logined {
		// redirect [GET]
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	// view
	resp := gin.H{
		"title":    "Accounts Page",
		"navi_key": "/accounts/",
	}
	ctx.HTML(http.StatusOK, "pages/accounts/accounts.tmpl", html.Response(resp, ctl.apiHeaderConf))
}
