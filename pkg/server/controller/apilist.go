package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// APILister interface
type APILister interface {
	APIListIndexAction(c *gin.Context)
	APIListIndex2Action(c *gin.Context)
}

// APIListIndexAction is top page for API List (react version)
func (ctl *controller) APIListIndexAction(c *gin.Context) {
	// return header and key
	ids, err := ctl.userRepo.GetUserIDs()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// View
	res := gin.H{
		"title":    "API List Page",
		"navi_key": "/apilist/",
		"ids":      ids,
	}
	c.HTML(http.StatusOK, "pages/apilist/index.tmpl", html.Response(res, ctl.apiHeader))
}

// APIListIndex2Action is top page for API List (this is old version)
func (ctl *controller) APIListIndex2Action(c *gin.Context) {
	ids, err := ctl.userRepo.GetUserIDs()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// View
	res := gin.H{
		"title":    "API List Page",
		"navi_key": "/apilist/",
		"ids":      ids,
	}
	c.HTML(http.StatusOK, "pages/apilist/index2.tmpl", html.Response(res, ctl.apiHeader))
}
