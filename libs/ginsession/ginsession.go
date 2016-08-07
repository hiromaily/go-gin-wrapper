package ginsession

import (
	"errors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	lg "github.com/hiromaily/golibs/log"
)

var sessionKey string = "secret123451234"

//TODO:setting to toml
func SetSession(r *gin.Engine, host, pass string) {

	var store sessions.RedisStore
	var err error
	if host == "" {
		//on memory
		store = sessions.NewCookieStore([]byte(sessionKey))
	} else {
		//session on Redis
		//store, err = sessions.NewRedisStore(80, "tcp", "localhost:6379", "", []byte("secret1234512345"))
		store, err = sessions.NewRedisStore(80, "tcp", host, pass, []byte(sessionKey))
		if err != nil {
			panic(err)
		}
	}

	strOptions := &sessions.Options{
		//Path: "/",
		//Domain: "/",   //It's better not to use.
		//MaxAge: 86400, //1day
		//MaxAge: 3600,  //1hour
		MaxAge:   300,   //5minutes
		Secure:   false, //TODO: set false in development environment, production environment requires true
		HttpOnly: true,
	}
	store.Options(*strOptions)
	r.Use(sessions.Sessions("ginsession", store))
}

//Set User
func SetUserSession(c *gin.Context, userId int) {
	session := sessions.Default(c)
	v := session.Get("uid")
	if v == nil {
		session.Set("uid", userId)
		session.Save()
	}
}

//Is Login
func IsLogin(c *gin.Context) (bRet bool, uid int) {
	session := sessions.Default(c)
	v := session.Get("uid")
	if v == nil {
		bRet = false
		uid = 0
	} else {
		bRet = true
		uid = v.(int)
	}
	//lg.Debugf("IsLogin: %v", bRet)
	return
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

//Set Token
func SetTokenSession(c *gin.Context, token string) {
	session := sessions.Default(c)
	session.Set("token", token)
	session.Save()
}

//Del Token
func DelTokenSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("token")
	session.Save()
}

//Get Token
func GetTokenSession(c *gin.Context) string {
	session := sessions.Default(c)
	v := session.Get("token")
	if v == nil {
		return ""
	}
	return v.(string)
}

//Check Token is valid
func IsTokenSessionValid(c *gin.Context, token string) bool {
	//default action
	if GetTokenSession(c) != token {
		//token error
		lg.Debug("Token is invalid.")

		//token delete
		DelTokenSession(c)

		c.AbortWithError(400, errors.New(("Token is invalid.")))
		return false
	}

	return true
}

//Set Count
func SetCountSession(c *gin.Context) {
	session := sessions.Default(c)
	var count int = 0
	v := session.Get("count")
	if v != nil {
		count = v.(int) + 1
	}
	session.Set("count", count)
	session.Save()
	lg.Debugf("session count:%d", count)
}
