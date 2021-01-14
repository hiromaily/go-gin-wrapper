package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	lg "github.com/hiromaily/golibs/log"
)

// ChatIndexAction is top page of chat [GET]
func (ctl *Controller) ChatIndexAction(c *gin.Context) {
	lg.Info("SignInGoogleAction()")

	// View
	c.HTML(http.StatusOK, "pages/chat/index.tmpl", gin.H{
		"title": "Chat Page",
	})
}
