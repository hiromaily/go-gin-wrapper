package users

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//Users
func UsersAction(c *gin.Context) {
	//Param

	//Logic

	//View
	c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
		"title": "Main website",
	})

}
