package bases

import (
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	"github.com/hiromaily/go-gin-wrapper/libs/csrf"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	hh "github.com/hiromaily/go-gin-wrapper/libs/httpheader"
	"github.com/hiromaily/go-gin-wrapper/libs/login"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
)

//TODO:define as common use.
func debugContext(c *gin.Context) {
	lg.Debugf("[c *gin.Context]: %#v \n", c)
	lg.Debugf("[c.Keys]: %#v \n", c.Keys)
	lg.Debugf("[c.Request.Method]: %s \n", c.Request.Method)
	lg.Debugf("[c.Request.Header]: %#v \n", c.Request.Header)
	lg.Debugf("[c.Request.Body]: %#v \n", c.Request.Body)
	lg.Debugf("[c.Request.URL]: %#v \n", c.Request.URL)
	lg.Debugf("[c.Value(ajax)]: %s \n", c.Value("ajax"))
	lg.Debugf("[hh.GetUrl(c)]: %s \n", hh.GetUrl(c))
	lg.Debugf("[hh.GetProto(c)]: %s \n", hh.GetProto(c))
}

// response for Login Page
func resLogin(c *gin.Context, input *login.LoginRequest, msg string, errors []string) {
	//token
	token := csrf.CreateToken()
	sess.SetTokenSession(c, token)

	//when crossing request, context data can't be left.
	//c.Set("getlogin", "xxx")

	if msg == "" {
		msg = "Enter Details to Login!!"
	}

	//it's better to not return nil
	if input == nil {
		input = &login.LoginRequest{}
	}

	//View
	c.HTML(http.StatusOK, "pages/bases/login.tmpl", gin.H{
		"message":  msg,
		"input":    input,
		"errors":   errors,
		"gintoken": token,
	})
}

//Index
func IndexAction(c *gin.Context) {
	//debug log
	//debugContext(c)

	//return header and key
	api := conf.GetConf().Api

	lg.Debugf("api.Header: %#v\n", api.Header)
	lg.Debugf("api.Key: %#v\n", api.Key)

	//View
	c.HTML(http.StatusOK, "pages/bases/index.tmpl", gin.H{
		"title":    "Top Page",
		"navi_key": "/",
		"header":   api.Header,
		"key":      api.Key,
	})
}

//Login [GET]
func LoginGetAction(c *gin.Context) {
	//debug log
	//debugContext(c)

	//If already loged in, go another page using redirect
	//Judge loged in or not.
	if bRet, id := sess.IsLogin(c); bRet {
		lg.Debugf("id: %d", id)

		//Redirect[GET]
		//FIXME:Browser request cache data when redirecting at status code 301
		//https://infra.xyz/archives/75
		//301 Moved Permanently   (Do cache,   it's possible to change from POST to GET)
		//302 Found               (Not cache,  it's possible to change from POST to GET)
		//307 Temporary Redirect  (Not cache,  it's not possible to change from POST to GET)
		//308 Moved Permanently   (Do cache,   it's not possible to change from POST to GET)

		//c.Redirect(http.StatusMovedPermanently, "/accounts/") //301
		c.Redirect(http.StatusTemporaryRedirect, "/accounts/") //307

		return
	}

	//return
	resLogin(c, nil, "", nil)
}

//Login [POST]
func LoginPostAction(c *gin.Context) {
	//debug log
	//debugContext(c)

	//check login
	userId, posted, errs := login.CheckLoginHTML(c)
	if errs != nil {
		resLogin(c, posted, "", errs)
		return
	}

	//When login is successful
	//Session
	sess.SetUserSession(c, userId)

	//token delete
	sess.DelTokenSession(c)

	//Change method POST to GET
	//Redirect[GET]
	//Status code 307 can't change post to get, 302 is suitable
	c.Redirect(http.StatusFound, "/accounts")

	return
}

//Logout [POST]
func LogoutPostAction(c *gin.Context) {
	lg.Debug("LogoutPostAction")
	//lg.Debug(sess.IsLogin(c))

	//Session
	sess.Logout(c)

	//lg.Debug(sess.IsLogin(c))
	api := conf.GetConf().Api

	//View
	c.HTML(http.StatusOK, "pages/bases/logout.tmpl", gin.H{
		"title":    "Logout Page",
		"navi_key": "/logout",
		"header":   api.Header,
		"key":      api.Key,
	})
}

//Logout [PUT] For Ajax
func LogoutPutAction(c *gin.Context) {
	lg.Debug("LogoutPutAction")
	//lg.Debug(sess.IsLogin(c))

	//Session
	sess.Logout(c)

	//lg.Debug(sess.IsLogin(c))

	//View
	c.JSON(http.StatusOK, gin.H{
		"message": "logout",
	})
}
