package admins

import (
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

// get user, it was set by the BasicAuth middleware

//Index [GET]
func IndexAction(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	lg.Debugf("[---]gin.AuthUserKey: %s", user)

	api := conf.GetConf().Api

	//View
	c.HTML(http.StatusOK, "pages/admins/gallery.tmpl", gin.H{
		"title":    "Admin Page",
		"navi_key": "/admin/",
		"header":   api.Header,
		"key":      api.Key,
	})
}
