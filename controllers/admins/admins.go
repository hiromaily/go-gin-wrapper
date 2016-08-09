package admins

import (
	"github.com/gin-gonic/gin"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

// get user, it was set by the BasicAuth middleware

//Index [GET]
func IndexAction(c *gin.Context) {
	//Param

	//Logic
	user := c.MustGet(gin.AuthUserKey).(string)
	lg.Debugf("user: %s", user)

	//View
	c.HTML(http.StatusOK, "pages/admins/gallery.tmpl", gin.H{
		"title": "Main website",
	})
}
