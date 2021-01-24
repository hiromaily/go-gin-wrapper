package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Errorer interface
type Errorer interface {
	Error404Action(ctx *gin.Context)
	Error405Action(ctx *gin.Context)
}

// Error404Action returns 404 error [GET]
// - WIP
func (ctl *controller) Error404Action(ctx *gin.Context) {
	referer := ctx.Request.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	// view
	ctx.HTML(http.StatusNotFound, "pages/errors/error.tmpl", gin.H{
		"message": "404 errors",
		"url":     referer,
	})
}

// Error405Action is for 405 error
// - WIP
func (ctl *controller) Error405Action(ctx *gin.Context) {
	referer := ctx.Request.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	// view
	ctx.HTML(http.StatusMethodNotAllowed, "pages/errors/error.tmpl", gin.H{
		"message": "405 errors",
		"url":     referer,
	})
}
