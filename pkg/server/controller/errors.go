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

// Error404Action is for 404 error [GET]
func (ctl *controller) Error404Action(ctx *gin.Context) {
	refURL := "/"
	if ctx.Request.Header.Get("Referer") != "" {
		refURL = ctx.Request.Header.Get("Referer")
	}

	ctx.HTML(http.StatusNotFound, "pages/errors/error.tmpl", gin.H{
		"code":    http.StatusNotFound,
		"message": "No where!",
		"url":     refURL,
	})

	// View
	ctx.HTML(http.StatusNotFound, "pages/errors/error.tmpl", gin.H{
		"message": "404 errors",
	})
}

// Error405Action is for 405 error
func (ctl *controller) Error405Action(ctx *gin.Context) {
	refURL := ctx.Request.Header.Get("Referer")

	// View
	ctx.HTML(http.StatusMethodNotAllowed, "pages/errors/error.tmpl", gin.H{
		"code":    http.StatusNotFound,
		"message": "405 errors",
		"url":     refURL,
	})
}
