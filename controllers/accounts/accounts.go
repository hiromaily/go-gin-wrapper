package accounts

import (
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	"github.com/hiromaily/go-gin-wrapper/libs/response/html"
	lg "github.com/hiromaily/golibs/log"
	//gin "gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-gonic/gin"
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
