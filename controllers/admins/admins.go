package admins

import (
	"github.com/hiromaily/go-gin-wrapper/libs/response/html"
	lg "github.com/hiromaily/golibs/log"
	//gin "gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

// get user, it was set by the BasicAuth middleware

//IndexAction [GET]
func IndexAction(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	lg.Debugf("[---]gin.AuthUserKey: %s", user)

	//View
	res := gin.H{
		"title":    "Admin Page",
		"navi_key": "/admin/",
	}
	c.HTML(http.StatusOK, "pages/admins/gallery.tmpl", html.Response(res))
}
