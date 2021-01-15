package ginsession

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
)

// SetSession is for session
func SetSession(r *gin.Engine, logger *zap.Logger, host, pass string, ses config.SessionConfig) {
	var store sessions.RedisStore
	var err error

	if host != "" {
		// session on Redis
		store, err = sessions.NewRedisStore(80, "tcp", host, pass, []byte(ses.Key))
		if err != nil {
			logger.Error("fail to call sessions.NewRedisStore()", zap.Error(err))
			// on memory
			store = sessions.NewCookieStore([]byte(ses.Key))
		}
	} else {
		// on memory
		store = sessions.NewCookieStore([]byte(ses.Key))
	}

	strOptions := &sessions.Options{
		// Path: "/",
		// Domain: "/",   //It's better not to use.
		// MaxAge: 86400, //1day
		// MaxAge: 3600,  //1hour
		MaxAge:   ses.MaxAge, // 5minutes
		Secure:   ses.Secure, // TODO: Set false in development environment, production environment requires true
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
	// logger.Debugf("IsLogin: %v", bRet)
	return
}

// Logout is to clear session
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

// SetTokenSession is to set token
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
func IsTokenSessionValid(c *gin.Context, logger *zap.Logger, token string) bool {

	logger.Info("IsTokenSessionValid",
		zap.String("GetTokenSession()", GetTokenSession(c)),
		zap.String("token", token),
	)

	var err error
	if GetTokenSession(c) == "" && token == "" {
		err = errors.New("Token is not allowed as blank.")
	} else if GetTokenSession(c) == "" {
		err = errors.New("Token is missing. Session might have expired.")
	} else if GetTokenSession(c) != token {
		err = errors.New("Token is invalid.")
	} else {
		return true
	}

	// token delete
	DelTokenSession(c)
	logger.Error("session error", zap.Error(err))
	c.AbortWithError(400, err)
	return false
}

// SetCountSession is for test
// TODO:delete this func
func SetCountSession(c *gin.Context, logger *zap.Logger) {
	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v != nil {
		count = v.(int) + 1
	}
	session.Set("count", count)
	session.Save()
	logger.Debug("SetCountSession", zap.Int("session count", count))
}
