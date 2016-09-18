package apilist

import (
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	models "github.com/hiromaily/go-gin-wrapper/models/mysql"
	//lg "github.com/hiromaily/golibs/log"
	"net/http"
)

// IndexAction is top page for API List (react version)
func IndexAction(c *gin.Context) {
	//return header and key
	api := conf.GetConf().Auth.API
	//lg.Debugf("api.Header: %#v\n", api.Header)
	//lg.Debugf("api.Key: %#v\n", api.Key)

	//Get User ids
	type UserID struct {
		ID int
	}
	var ids []UserID

	err := models.GetDB().GetUserIds(&ids)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	//View
	c.HTML(http.StatusOK, "pages/apilist/index.tmpl", gin.H{
		"title":    "API List Page",
		"navi_key": "/apilist/",
		"ids":      ids,
		"header":   api.Header,
		"key":      api.Key,
	})
}

// Index2Action is top page for API List (this is old version)
func Index2Action(c *gin.Context) {
	//debug log
	//debugContext(c)

	//return header and key
	api := conf.GetConf().Auth.API
	//lg.Debugf("api.Header: %#v\n", api.Header)
	//lg.Debugf("api.Key: %#v\n", api.Key)

	//Get User ids
	type UserID struct {
		ID int
	}
	var ids []UserID

	err := models.GetDB().GetUserIds(&ids)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	//View
	c.HTML(http.StatusOK, "pages/apilist/index2.tmpl", gin.H{
		"title":    "API List Page",
		"navi_key": "/apilist/",
		"ids":      ids,
		"header":   api.Header,
		"key":      api.Key,
	})
}
