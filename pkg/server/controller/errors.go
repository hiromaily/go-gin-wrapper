package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Errorer interface
type Errorer interface {
	Error404Action(c *gin.Context)
	Error405Action(c *gin.Context)
}

// Error404Action is for 404 error [GET]
func (ctl *controller) Error404Action(c *gin.Context) {
	refURL := "/"
	if c.Request.Header.Get("Referer") != "" {
		refURL = c.Request.Header.Get("Referer")
	}

	c.HTML(http.StatusNotFound, "pages/errors/error.tmpl", gin.H{
		"code":    http.StatusNotFound,
		"message": "No where!",
		"url":     refURL,
	})

	// View
	c.HTML(http.StatusNotFound, "pages/errors/error.tmpl", gin.H{
		"message": "404 errors",
	})
}

// Error405Action is for 405 error
func (ctl *controller) Error405Action(c *gin.Context) {
	refURL := c.Request.Header.Get("Referer")

	// View
	c.HTML(http.StatusMethodNotAllowed, "pages/errors/error.tmpl", gin.H{
		"code":    http.StatusNotFound,
		"message": "405 errors",
		"url":     refURL,
	})
}
