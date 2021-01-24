package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// APILister interface
type APILister interface {
	APIListIndexAction(ctx *gin.Context)
	APIListIndex2Action(ctx *gin.Context)
}

// APIListIndexAction returns API List page
func (ctl *controller) APIListIndexAction(ctx *gin.Context) {
	ids, err := ctl.userRepo.GetUserIDs()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// view
	res := gin.H{
		"title":    "API List Page",
		"navi_key": "/apilist/",
		"ids":      ids,
	}
	// return header and key as well
	ctx.HTML(http.StatusOK, "pages/apilist/index.tmpl", html.Response(res, ctl.apiHeaderConf))
}

// APIListIndex2Action returns API List page (old version)
func (ctl *controller) APIListIndex2Action(ctx *gin.Context) {
	ids, err := ctl.userRepo.GetUserIDs()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// View
	res := gin.H{
		"title":    "API List Page",
		"navi_key": "/apilist/",
		"ids":      ids,
	}
	ctx.HTML(http.StatusOK, "pages/apilist/index2.tmpl", html.Response(res, ctl.apiHeaderConf))
}
