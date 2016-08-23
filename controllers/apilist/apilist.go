package apilist

import (
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	models "github.com/hiromaily/go-gin-wrapper/models/mysql"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

//Index
func IndexAction(c *gin.Context) {
	//debug log
	//debugContext(c)

	//return header and key
	api := conf.GetConf().Api
	lg.Debugf("api.Header: %#v\n", api.Header)
	lg.Debugf("api.Key: %#v\n", api.Key)

	//Get User ids
	type UserId struct {
		Id int
	}
	var ids []UserId

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
