package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/ginctx"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// Baser interface
type Baser interface {
	BaseIndexAction(ctx *gin.Context)
	BaseLoginGetAction(ctx *gin.Context)
	BaseLoginPostAction(ctx *gin.Context)
	BaseLogoutPostAction(ctx *gin.Context)
	BaseLogoutPutAction(ctx *gin.Context)
}

// BaseIndexAction returns top page
func (ctl *controller) BaseIndexAction(ctx *gin.Context) {
	ctl.logger.Debug("BaseIndexAction")

	// debug log
	ginctx.DebugContext(ctx, ctl.logger)

	// view
	res := gin.H{
		"title":    "Top Page",
		"navi_key": "/",
	}
	ctx.HTML(http.StatusOK, "pages/bases/index.tmpl", html.Response(res, ctl.apiHeaderConf))
}

// BaseLoginGetAction returns login page [GET]
func (ctl *controller) BaseLoginGetAction(ctx *gin.Context) {
	ctl.logger.Debug("BaseLoginGetAction")

	// debug log
	// ctl.debugContext(ctx)

	if logined, _ := ctl.session.IsLogin(ctx); logined {
		// Redirect[GET]
		// Note: Browser uses cache data when redirecting with status code 301
		// 301 Moved Permanently   (Do cache,   it's possible to change from POST to GET)
		// 302 Found               (Not cache,  it's possible to change from POST to GET)
		// 307 Temporary Redirect  (Not cache,  it's not possible to change from POST to GET)
		// 308 Moved Permanently   (Do cache,   it's not possible to change from POST to GET)
		ctx.Redirect(http.StatusTemporaryRedirect, "/accounts/") // 307
		return
	}

	ctl.loginResponse(ctx, nil, "", nil)
}

// BaseLoginPostAction handles login request [POST]
func (ctl *controller) BaseLoginPostAction(ctx *gin.Context) {
	ctl.logger.Debug("BaseLoginPostAction")

	userID, loginRequest, errs := ctl.login(ctx)
	if len(errs) != 0 {
		ctl.logger.Debug("login_error", zap.Any("errors", errs))
		ctl.loginResponse(ctx, loginRequest, "", errs)
		return
	}

	// start session
	ctl.session.SetUserID(ctx, userID)
	ctl.session.DeleteToken(ctx)

	// change method POST to GET
	// Redirect[GET]
	// Status code 307 can't change post to get, 302 is suitable
	ctx.Redirect(http.StatusFound, "/accounts/")
}

// BaseLogoutPostAction handles logout [POST]
func (ctl *controller) BaseLogoutPostAction(ctx *gin.Context) {
	ctl.logger.Debug("LogoutPostAction")

	// Session
	ctl.session.Logout(ctx)

	// View
	res := gin.H{
		"title":    "Logout Page",
		"navi_key": "/logout",
	}
	ctx.HTML(http.StatusOK, "pages/bases/logout.tmpl", html.Response(res, ctl.apiHeaderConf))
}

// BaseLogoutPutAction handles logout API [PUT]
func (ctl *controller) BaseLogoutPutAction(ctx *gin.Context) {
	ctl.logger.Debug("LogoutPutAction")

	ctl.session.Logout(ctx)

	// view
	ctx.JSON(http.StatusOK, gin.H{
		"message": "logout",
	})
}
