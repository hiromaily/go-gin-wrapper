package accounts

import (
	"github.com/gin-gonic/gin"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	"net/http"
)

//Accounts [GET]
func AccountsGetAction(c *gin.Context) {
	//judge login
	if bRet, _ := sess.IsLogin(c); !bRet {
		//Redirect[GET]
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	//View
	c.HTML(http.StatusOK, "pages/accounts/accounts.tmpl", gin.H{
		"title":    "Accounts Page",
		"navi_key": "/accounts/",
	})
}
