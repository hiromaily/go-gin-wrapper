package news

import (
	"github.com/hiromaily/go-gin-wrapper/core/libs/response/html"
	models "github.com/hiromaily/go-gin-wrapper/core/models/mongo"
	//gin "gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ParamNews is for news data from MongoDB
type ParamNews struct {
	Classes  []string
	Articles []models.Articles
}

// IndexAction is for top page of news [GET]
func IndexAction(c *gin.Context) {
	//Get news
	articles, err := models.GetDB().GetArticlesData(0)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	//Param
	className := []string{"alert-success", "alert-info", "alert-warning", "alert-danger"}

	//View
	res := gin.H{
		"title":    "News Page",
		"navi_key": "/news/",
		"articles": articles,
		"class":    className,
	}
	c.HTML(http.StatusOK, "pages/news/news.tmpl", html.Response(res))
}

// Index2Action is still working in progress.
func Index2Action(c *gin.Context) {
	//Get news
	items, err := models.GetDB().GetArticlesData2(0)
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
