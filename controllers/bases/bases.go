package bases

import (
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	"github.com/hiromaily/go-gin-wrapper/libs/csrf"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	hh "github.com/hiromaily/go-gin-wrapper/libs/httpheader"
	"github.com/hiromaily/go-gin-wrapper/models"
	lg "github.com/hiromaily/golibs/log"
	valid "github.com/hiromaily/golibs/validator"
	"net/http"
)

type LoginRequest struct {
	Email string `valid:"nonempty,email,min=5,max=40" field:"email" dispName:"E-Mail"`
	Pass  string `valid:"nonempty,min=8,max=20" field:"pass" dispName:"Password"`
}

//TODO: this should be defined as something common library.
var ErrFmt = map[string]string{
	"nonempty": "Empty is not allowed on %s",
	"email":    "Format of %s is invalid",
	"alphanum": "Only alphabet is allowd on %s",
	"min":      "At least %s of characters is required on %s",
	"max":      "At a maximum %s of characters is allowed on %s",
}

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

//Index
func IndexAction(c *gin.Context) {
	//debug log
	debugContext(c)

	//return header and key
	api := conf.GetConfInstance().Api

	lg.Debugf("api.Header: %#v\n", api.Header)
	lg.Debugf("api.Key: %#v\n", api.Key)

	//View
	c.HTML(http.StatusOK, "pages/bases/index.tmpl", gin.H{
		"title":  "Main website",
		"header": api.Header,
		"key":    api.Key,
	})
}

//Login [GET]
func LoginGetAction(c *gin.Context) {
	//debug log
	debugContext(c)

	//If already loged in, go another page using redirect
	//Judge loged in or not.
	if bRet, id := sess.IsLogin(c); bRet {
		lg.Debugf("id: %d", id)

		//Redirect[GET]
		//FIXME:Browser cache request when redirecting at status code 301
		//https://infra.xyz/archives/75
		//301 Moved Permanently   (Do cache,   it's possible to change from POST to GET)
		//302 Found               (Not cache,  it's possible to change from POST to GET)
		//307 Temporary Redirect  (Not cache,  it's not possible to change from POST to GET)
		//308 Moved Permanently   (Do cache,   it's not possible to change from POST to GET)

		//c.Redirect(http.StatusMovedPermanently, "/accounts") //301
		c.Redirect(http.StatusTemporaryRedirect, "/accounts") //307

		return
	}

	//token
	token := csrf.CreateToken()
	sess.SetTokenSession(c, token)

	//when crossing request, context data can't be left.
	//c.Set("getlogin", "xxx")

	//View
	c.HTML(http.StatusOK, "pages/bases/login.tmpl", gin.H{
		"message":  "nothing special",
		"gintoken": token,
	})
}

//Login [POST]
func LoginPostAction(c *gin.Context) {
	//debug log
	debugContext(c)

	//Get Post Parameters
	inputEmail := c.PostForm("inputEmail") //return is string type
	inputPassword := c.PostForm("inputPassword")
	//tokenPosted := c.PostForm("gintoken")
	lg.Debugf("inputEmail : %s\n", inputEmail)
	lg.Debugf("inputPassword : %s\n", inputPassword)
	//lg.Debugf("gintoken : %s\n", tokenPosted)

	//Validation
	posted := &LoginRequest{Email: inputEmail, Pass: inputPassword}
	//FIXME: It doesn't work when passed address of struct type.
	mRet := valid.CheckValidation(posted, false)
	//map[string][]string{"pass":[]string{"min"}, "test":[]string{"nonempty"}}
	if len(mRet) != 0 {
		msgs := valid.ConvertErrorMsg(mRet, ErrFmt)

		lg.Debugf("validation error: %#v", msgs)

		//token
		token := csrf.CreateToken()
		sess.SetTokenSession(c, token)

		//View
		c.HTML(http.StatusOK, "pages/bases/login.tmpl", gin.H{
			"message":  "validation error happend",
			"gintoken": token,
		})
		return
	}

	//Check inputed mail and password
	//aaaa@aa.jp / password
	userId, err := models.GetDBInstance().IsUserEmail(inputEmail, inputPassword)
	if err != nil {
		lg.Debugf("login error: %#v", mRet)

		//token
		token := csrf.CreateToken()
		sess.SetTokenSession(c, token)

		//View
		c.HTML(http.StatusOK, "pages/bases/login.tmpl", gin.H{
			"message":  "mailaddress and password may be wrong",
			"gintoken": token,
		})
		return
	}

	//Session
	sess.SetUserSession(c, userId)

	//token delete
	sess.DelTokenSession(c)

	//Change method post to get
	//Redirect[GET]
	//c.Redirect(http.StatusMovedPermanently, "/accounts")
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

	//View
	c.HTML(http.StatusOK, "pages/bases/logout.tmpl", gin.H{
		"message": "logout was done.",
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
