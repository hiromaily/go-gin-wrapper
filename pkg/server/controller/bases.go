package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/csrf"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
)

// Baser interface
type Baser interface {
	BaseIndexAction(c *gin.Context)
	BaseLoginGetAction(c *gin.Context)
	BaseLoginPostAction(c *gin.Context)
	BaseLogoutPostAction(c *gin.Context)
	BaseLogoutPutAction(c *gin.Context)
}

// nolint: unused, deadcode
func debugContext(c *gin.Context, logger *zap.Logger) {
	logger.Debug("request",
		zap.Any("gin_context", c),
		zap.Any("gin_keys", c.Keys),
		zap.String("request_method", c.Request.Method),
		zap.Any("request_header", c.Request.Header),
		zap.Any("request_body", c.Request.Body),
		zap.Any("request_url", c.Request.URL),
		zap.Any("request_ajax", c.Value("ajax")),
		zap.String("request_get_url", hh.GetURL(c)),
		zap.String("request_get_protocol", hh.GetProto(c)),
	)
}

// response for Login Page
func (ctl *controller) resLogin(c *gin.Context, input *LoginRequest, msg string, errors []string) {
	// token
	token := csrf.CreateToken(ctl.logger)
	sess.SetTokenSession(c, token)

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
	c.HTML(http.StatusOK, "pages/bases/login.tmpl", gin.H{
		"message":               msg,
		"input":                 input,
		"github.com/pkg/errors": errors,
		"gintoken":              token,
		"gURL":                  gURL,
		"fURL":                  fURL,
	})
}

// BaseIndexAction is top page
func (ctl *controller) BaseIndexAction(c *gin.Context) {
	// debug log
	// debugContext(c)

	// View
	res := gin.H{
		"title":    "Top Page",
		"navi_key": "/",
	}
	c.HTML(http.StatusOK, "pages/bases/index.tmpl", html.Response(res, ctl.apiHeader))
}

// BaseLoginGetAction is for login page [GET]
func (ctl *controller) BaseLoginGetAction(c *gin.Context) {
	// debug log
	// debugContext(c)

	// If already loged in, go another page using redirect
	// Judge loged in or not.
	if bRet, _ := sess.IsLogin(c); bRet {
		// Redirect[GET]
		// FIXME: Browser request cache data when redirecting at status code 301
		// https://infra.xyz/archives/75
		// 301 Moved Permanently   (Do cache,   it's possible to change from POST to GET)
		// 302 Found               (Not cache,  it's possible to change from POST to GET)
		// 307 Temporary Redirect  (Not cache,  it's not possible to change from POST to GET)
		// 308 Moved Permanently   (Do cache,   it's not possible to change from POST to GET)

		// c.Redirect(http.StatusMovedPermanently, "/accounts/") //301
		c.Redirect(http.StatusTemporaryRedirect, "/accounts/") // 307

		return
	}

	// return
	ctl.resLogin(c, nil, "", nil)
}

// BaseLoginPostAction is to receive user request from login page [POST]
func (ctl *controller) BaseLoginPostAction(c *gin.Context) {
	// check login
	userID, posted, errs := ctl.CheckLoginOnHTML(c)
	if errs != nil {
		ctl.resLogin(c, posted, "", errs)
		return
	}

	// When login is successful
	// Session
	sess.SetUserSession(c, userID)

	// token delete
	sess.DelTokenSession(c)

	// Change method POST to GET
	// Redirect[GET]
	// Status code 307 can't change post to get, 302 is suitable
	c.Redirect(http.StatusFound, "/accounts/")
}

// BaseLogoutPostAction is for logout [POST]
func (ctl *controller) BaseLogoutPostAction(c *gin.Context) {
	ctl.logger.Debug("LogoutPostAction")

	// Session
	sess.Logout(c)

	// View
	res := gin.H{
		"title":    "Logout Page",
		"navi_key": "/logout",
	}
	c.HTML(http.StatusOK, "pages/bases/logout.tmpl", html.Response(res, ctl.apiHeader))
}

// BaseLogoutPutAction is for logout by Ajax [PUT]
func (ctl *controller) BaseLogoutPutAction(c *gin.Context) {
	ctl.logger.Debug("LogoutPutAction")

	// Session
	sess.Logout(c)

	// View
	c.JSON(http.StatusOK, gin.H{
		"message": "logout",
	})
}
