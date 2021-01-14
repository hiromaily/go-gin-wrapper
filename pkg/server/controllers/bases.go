package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/csrf"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/response/html"
	lg "github.com/hiromaily/golibs/log"
)

// TODO:define as common use.
// nolint: unused, deadcode
func debugContext(c *gin.Context) {
	lg.Debugf("[c *gin.Context]: %#v \n", c)
	lg.Debugf("[c.Keys]: %#v \n", c.Keys)
	lg.Debugf("[c.Request.Method]: %s \n", c.Request.Method)
	lg.Debugf("[c.Request.Header]: %#v \n", c.Request.Header)
	lg.Debugf("[c.Request.Body]: %#v \n", c.Request.Body)
	lg.Debugf("[c.Request.URL]: %#v \n", c.Request.URL)
	lg.Debugf("[c.Value(ajax)]: %s \n", c.Value("ajax"))
	lg.Debugf("[hh.GetUrl(c)]: %s \n", hh.GetURL(c))
	lg.Debugf("[hh.GetProto(c)]: %s \n", hh.GetProto(c))
}

// response for Login Page
func (ctl *Controller) resLogin(c *gin.Context, input *LoginRequest, msg string, errors []string) {
	// token
	token := csrf.CreateToken()
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
func (ctl *Controller) BaseIndexAction(c *gin.Context) {
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
func (ctl *Controller) BaseLoginGetAction(c *gin.Context) {
	// debug log
	// debugContext(c)

	// If already loged in, go another page using redirect
	// Judge loged in or not.
	if bRet, id := sess.IsLogin(c); bRet {
		lg.Debugf("id: %d", id)

		// Redirect[GET]
		// FIXME:Browser request cache data when redirecting at status code 301
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
func (ctl *Controller) BaseLoginPostAction(c *gin.Context) {
	// debug log
	// debugContext(c)

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
func (ctl *Controller) BaseLogoutPostAction(c *gin.Context) {
	lg.Debug("LogoutPostAction")
	// lg.Debug(sess.IsLogin(c))

	// Session
	sess.Logout(c)

	// View
	// View
	res := gin.H{
		"title":    "Logout Page",
		"navi_key": "/logout",
	}
	c.HTML(http.StatusOK, "pages/bases/logout.tmpl", html.Response(res, ctl.apiHeader))
}

// BaseLogoutPutAction is for logout by Ajax [PUT]
func (ctl *Controller) BaseLogoutPutAction(c *gin.Context) {
	lg.Debug("LogoutPutAction")
	// lg.Debug(sess.IsLogin(c))

	// Session
	sess.Logout(c)

	// lg.Debug(sess.IsLogin(c))

	// View
	c.JSON(http.StatusOK, gin.H{
		"message": "logout",
	})
}
