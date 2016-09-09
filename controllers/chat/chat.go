package chat

import (
	"github.com/gin-gonic/gin"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

// IndexAction is top page of chat [GET]
func IndexAction(c *gin.Context) {
	lg.Info("SignInGoogleAction()")

	//View
	c.HTML(http.StatusOK, "pages/chat/index.tmpl", gin.H{
		"title": "Chat Page",
	})
}
