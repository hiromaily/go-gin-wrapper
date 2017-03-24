package ginsession

import (
	"errors"
	"github.com/gin-gonic/contrib/sessions"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	//"github.com/hiromaily/go-gin-wrapper/libs/sessions"
	lg "github.com/hiromaily/golibs/log"
	//gin "gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-gonic/gin"
)

// SetSession is for session
func SetSession(r *gin.Engine, host, pass string) {

	var store sessions.RedisStore
	var err error
	ses := conf.GetConf().Server.Session

	if host == "" {
		//on memory
		store = sessions.NewCookieStore([]byte(ses.Key))
	} else {
		//session on Redis
		//store, err = sessions.NewRedisStore(80, "tcp", "localhost:6379", "", []byte("secret1234512345"))
		store, err = sessions.NewRedisStore(80, "tcp", host, pass, []byte(ses.Key))
		if err != nil {
			panic(err)
		}
	}

	strOptions := &sessions.Options{
		//Path: "/",
		//Domain: "/",   //It's better not to use.
		//MaxAge: 86400, //1day
		//MaxAge: 3600,  //1hour
		MaxAge:   ses.MaxAge, //5minutes
		Secure:   ses.Secure, //TODO: Set false in development environment, production environment requires true
		HttpOnly: ses.HTTPOnly,
	}
	store.Options(*strOptions)
	r.Use(sessions.Sessions(ses.Name, store))
}

// SetUserSession is set user session data
func SetUserSession(c *gin.Context, userID int) {
	session := sessions.Default(c)
	v := session.Get("uid")
	if v == nil {
		session.Set("uid", userID)
		session.Save()
	}
}

// IsLogin is whether user have already loged in or not.
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

// Logout is to clear session
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

//SetTokenSession is to set token
func SetTokenSession(c *gin.Context, token string) {
	session := sessions.Default(c)
	session.Set("token", token)
	session.Save()
}

// DelTokenSession is to delete token
func DelTokenSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("token")
	session.Save()
}

// GetTokenSession is to get token
func GetTokenSession(c *gin.Context) string {
	session := sessions.Default(c)
	v := session.Get("token")
	if v == nil {
		return ""
	}
	return v.(string)
}

// IsTokenSessionValid is whether check token is valid or not
func IsTokenSessionValid(c *gin.Context, token string) bool {
	strErr := ""

	lg.Info("[IsTokenSessionValid]")
	lg.Debugf("GetTokenSession(c): %v", GetTokenSession(c))
	lg.Debugf("token: %v", token)

	if GetTokenSession(c) == "" && token == "" {
		strErr = "Token is not allowed as blank."
	} else if GetTokenSession(c) == "" {
		strErr = "Token is missing. Session might have expired."
	} else if GetTokenSession(c) != token {
		strErr = "Token is invalid."
	} else {
		return true
	}

	//token delete
	DelTokenSession(c)

	lg.Error(strErr)
	c.AbortWithError(400, errors.New((strErr)))
	return false
}

// SetCountSession is for test (TODO:delete this func)
func SetCountSession(c *gin.Context) {
	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v != nil {
		count = v.(int) + 1
	}
	session.Set("count", count)
	session.Save()
	lg.Debugf("session count:%d", count)
}
