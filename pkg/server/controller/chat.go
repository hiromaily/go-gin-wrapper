package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Chater interface
type Chater interface {
	ChatIndexAction(ctx *gin.Context)
}

// ChatIndexAction returns chat page [GET]
// - WIP
func (ctl *controller) ChatIndexAction(ctx *gin.Context) {
	ctl.logger.Info("ChatIndexAction")

	// View
	ctx.HTML(http.StatusOK, "pages/chat/index.tmpl", gin.H{
		"title": "Chat Page",
	})
}
