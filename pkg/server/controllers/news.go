package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/go-gin-wrapper/pkg/libs/response/html"
	models "github.com/hiromaily/go-gin-wrapper/pkg/models/mongo"
)

// ParamNews is for news data from MongoDB
type ParamNews struct {
	Classes  []string
	Articles []models.Articles
}

// NewsIndexAction is for top page of news [GET]
func (ctl *Controller) NewsIndexAction(c *gin.Context) {
	//Get news
	articles, err := ctl.mongo.GetArticlesData(0)
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

// NewsIndex2Action is still working in progress.
func (ctl *Controller) NewsIndex2Action(c *gin.Context) {
	//Get news
	items, err := ctl.mongo.GetArticlesData2(0)
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