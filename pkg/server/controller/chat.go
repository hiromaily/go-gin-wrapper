package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ChatIndexAction is top page of chat [GET]
func (ctl *Controller) ChatIndexAction(c *gin.Context) {
	ctl.logger.Info("ChatIndexAction")

	// View
	c.HTML(http.StatusOK, "pages/chat/index.tmpl", gin.H{
		"title": "Chat Page",
	})
}
