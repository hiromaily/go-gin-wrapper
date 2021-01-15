package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// AccountIndexAction [GET]
func (ctl *Controller) AccountIndexAction(c *gin.Context) {
	ctl.logger.Info("AccountIndexAction")

	// judge login
	if bRet, _ := sess.IsLogin(c); !bRet {
		// Redirect[GET]
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	// View
	res := gin.H{
		"title":    "Accounts Page",
		"navi_key": "/accounts/",
	}
	c.HTML(http.StatusOK, "pages/accounts/accounts.tmpl", html.Response(res, ctl.apiHeader))
}
