package errors

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//404 Error
func Error404Action(c *gin.Context) {
	//Param

	//View
	c.HTML(http.StatusOK, "errors/error.tmpl", gin.H{
		"message": "404 errors",
	})
}

//405 Error
func Error405Action(c *gin.Context) {
	//Param

	//View
	c.HTML(http.StatusOK, "errors/error.tmpl", gin.H{
		"message": "405 errors",
	})
}
