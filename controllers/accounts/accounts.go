package accounts

import (
	"github.com/gin-gonic/gin"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	"github.com/hiromaily/go-gin-wrapper/libs/response/html"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

//IndexAction [GET]
func IndexAction(c *gin.Context) {
	lg.Info("AccountsGetAction()")

	//judge login
	if bRet, _ := sess.IsLogin(c); !bRet {
		//Redirect[GET]
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	//View
	res := gin.H{
		"title":    "Accounts Page",
		"navi_key": "/accounts/",
	}
	c.HTML(http.StatusOK, "pages/accounts/accounts.tmpl", html.Response(res))
}
