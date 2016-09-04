package accounts

import (
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

//Accounts [GET]
func AccountsGetAction(c *gin.Context) {
	lg.Info("AccountsGetAction()")

	//judge login
	if bRet, _ := sess.IsLogin(c); !bRet {
		//Redirect[GET]
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	api := conf.GetConf().Auth.Api

	//View
	c.HTML(http.StatusOK, "pages/accounts/accounts.tmpl", gin.H{
		"title":    "Accounts Page",
		"navi_key": "/accounts/",
		"header":   api.Header,
		"key":      api.Key,
	})
}
