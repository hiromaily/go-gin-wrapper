package errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//404 Error [GET]
func Error404Action(c *gin.Context) {

	refUrl := "/"
	if c.Request.Header.Get("Referer") != "" {
		refUrl = c.Request.Header.Get("Referer")
	}

	c.HTML(http.StatusNotFound, "pages/errors/error.tmpl", gin.H{
		"code":    http.StatusNotFound,
		"message": "No where!",
		"url":     refUrl,
	})

	//View
	c.HTML(http.StatusNotFound, "pages/errors/error.tmpl", gin.H{
		"message": "404 errors",
	})
}

//405 Error
func Error405Action(c *gin.Context) {

	refUrl := c.Request.Header.Get("Referer")

	//View
	c.HTML(http.StatusMethodNotAllowed, "pages/errors/error.tmpl", gin.H{
		"code":    http.StatusNotFound,
		"message": "405 errors",
		"url":     refUrl,
	})
}
