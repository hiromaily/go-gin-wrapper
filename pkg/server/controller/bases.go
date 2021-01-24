package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/ginurl"
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

// nolint: unused, deadcode
func debugContext(ctx *gin.Context, logger *zap.Logger) {
	logger.Debug("request",
		zap.Any("gin_ctx", ctx),
		zap.Any("gin_ctx_keys", ctx.Keys),
		zap.String("request_method", ctx.Request.Method),
		zap.Any("request_header", ctx.Request.Header),
		zap.Any("request_body", ctx.Request.Body),
		zap.Any("request_url", ctx.Request.URL),
		zap.String("request_url_string", ginurl.GetURLString(ctx)),
		zap.String("request_protocol", ctx.Request.Proto),
		zap.Any("ctx_value_ajax", ctx.Value("ajax")),
	)
}

// response for Login Page
func (ctl *controller) resLogin(ctx *gin.Context, input *LoginRequest, msg string, errors []string) {
	// token
	token := ctl.session.GenerateToken()
	ctl.session.SetToken(ctx, token)

	// Google Open ID
	gURL := "/oauth2/google/signin"
	fURL := "/oauth2/facebook/signin"

	if msg == "" {
		msg = "Enter Details to Login!!"
	}

	// it's better to not return nil
	if input == nil {
		input = &LoginRequest{}
	}

	// View
	ctx.HTML(http.StatusOK, "pages/bases/login.tmpl", gin.H{
		"message":               msg,
		"input":                 input,
		"github.com/pkg/errors": errors,
		"gintoken":              token,
		"gURL":                  gURL,
		"fURL":                  fURL,
	})
}

// BaseIndexAction is top page
func (ctl *controller) BaseIndexAction(ctx *gin.Context) {
	// debug log
	// debugContext(ctx)

	// View
	res := gin.H{
		"title":    "Top Page",
		"navi_key": "/",
	}
	ctx.HTML(http.StatusOK, "pages/bases/index.tmpl", html.Response(res, ctl.apiHeaderConf))
}

// BaseLoginGetAction is for login page [GET]
func (ctl *controller) BaseLoginGetAction(ctx *gin.Context) {
	// debug log
	// debugContext(ctx)

	// If already loged in, go another page using redirect
	// Judge loged in or not.
	if bRet, _ := ctl.session.IsLogin(ctx); bRet {
		// Redirect[GET]
		// FIXME: Browser request cache data when redirecting at status code 301
		// https://infra.xyz/archives/75
		// 301 Moved Permanently   (Do cache,   it's possible to change from POST to GET)
		// 302 Found               (Not cache,  it's possible to change from POST to GET)
		// 307 Temporary Redirect  (Not cache,  it's not possible to change from POST to GET)
		// 308 Moved Permanently   (Do cache,   it's not possible to change from POST to GET)

		// ctx.Redirect(http.StatusMovedPermanently, "/accounts/") //301
		ctx.Redirect(http.StatusTemporaryRedirect, "/accounts/") // 307

		return
	}

	// return
	ctl.resLogin(ctx, nil, "", nil)
}

// BaseLoginPostAction is to receive user request from login page [POST]
func (ctl *controller) BaseLoginPostAction(ctx *gin.Context) {
	// check login
	userID, posted, errs := ctl.CheckLoginOnHTML(ctx)
	if errs != nil {
		ctl.resLogin(ctx, posted, "", errs)
		return
	}

	// When login is successful
	// Session
	ctl.session.SetUserID(ctx, userID)

	// token delete
	ctl.session.DeleteToken(ctx)

	// Change method POST to GET
	// Redirect[GET]
	// Status code 307 can't change post to get, 302 is suitable
	ctx.Redirect(http.StatusFound, "/accounts/")
}

// BaseLogoutPostAction is for logout [POST]
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

// BaseLogoutPutAction is for logout by Ajax [PUT]
func (ctl *controller) BaseLogoutPutAction(ctx *gin.Context) {
	ctl.logger.Debug("LogoutPutAction")

	// Session
	ctl.session.Logout(ctx)

	// View
	ctx.JSON(http.StatusOK, gin.H{
		"message": "logout",
	})
}
