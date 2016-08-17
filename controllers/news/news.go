package news

import (
	"github.com/gin-gonic/gin"
	models "github.com/hiromaily/go-gin-wrapper/models/mongo"
	"net/http"
)

//News [GET]
func NewsGetAction(c *gin.Context) {
	//Get news
	items, err := models.GetDB().GetArticlesData(0)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	//View
	c.HTML(http.StatusOK, "pages/news/news.tmpl", gin.H{
		"title":    "News Page",
		"navi_key": "/news",
		"items":    items,
	})
}
