package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
	lg "github.com/hiromaily/golibs/log"
)

// get user, it was set by the BasicAuth middleware

//AdminIndexAction [GET]
func (ctl *Controller) AdminIndexAction(c *gin.Context) {
	user := c.MustGet(gin.AuthUserKey).(string)
	lg.Debugf("[---]gin.AuthUserKey: %s", user)

	//View
	res := gin.H{
		"title":    "Admin Page",
		"navi_key": "/admin/",
	}
	c.HTML(http.StatusOK, "pages/admins/gallery.tmpl", html.Response(res))
}
