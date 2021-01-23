package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// Accounter interface
type Accounter interface {
	AccountIndexAction(c *gin.Context)
}

// AccountIndexAction [GET]
func (ctl *controller) AccountIndexAction(c *gin.Context) {
	ctl.logger.Info("controller AccountIndexAction")

	// validate access
	if logined, _ := sess.IsLogin(c); !logined {
		// redirect [GET]
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	// view
	resp := gin.H{
		"title":    "Accounts Page",
		"navi_key": "/accounts/",
	}
	c.HTML(http.StatusOK, "pages/accounts/accounts.tmpl", html.Response(resp, ctl.apiHeader))
}
