package news

import (
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	models "github.com/hiromaily/go-gin-wrapper/models/mongo"
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
	//params := ParamNews{
	//	Classes: []string{"alert-success","alert-info","alert-warning","alert-danger"},
	//	Articles: articles,
	//}
	className := []string{"alert-success", "alert-info", "alert-warning", "alert-danger"}

	api := conf.GetConf().Auth.API

	//View
	c.HTML(http.StatusOK, "pages/news/news.tmpl", gin.H{
		"title":    "News Page",
		"navi_key": "/news/",
		//"params": params,
		"articles": articles,
		"class":    className,
		"header":   api.Header,
		"key":      api.Key,
	})
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
