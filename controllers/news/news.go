package news

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//News [GET]
func NewsGetAction(c *gin.Context) {
	//Param

	//Logic

	//only session
	//_ = sess.IsLogin(c)

	//View
	c.HTML(http.StatusOK, "pages/news/news.tmpl", gin.H{
		"title":    "News Page",
		"navi_key": "/news",
	})

}
