package ginsession

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
)

// NewRedisStore returns redis session store
func NewRedisStore(logger *zap.Logger, host, pass string, conf *config.Session) sessions.Store {
	store, err := sessions.NewRedisStore(80, "tcp", host, pass, []byte(conf.Key))
	if err != nil {
		logger.Error("fail to call sessions.NewRedisStore()", zap.Error(err))
		// on memory
		store = sessions.NewCookieStore([]byte(conf.Key))
	}
	return store
}

// NewCookieStore returns cookie store
func NewCookieStore(conf *config.Session) sessions.Store {
	return sessions.NewCookieStore([]byte(conf.Key))
}

// SetOption sets option to session store
func SetOption(store sessions.RedisStore, conf *config.Session) {
	strOptions := &sessions.Options{
		// Path: "/",
		// Domain: "/",   // It's better not to use
		// MaxAge: 86400, // 1 day
		// MaxAge: 3600,  // 1 hour
		MaxAge:   conf.MaxAge, // 5 minutes
		Secure:   conf.Secure, // TODO: Set false in development environment, production environment requires true
		HttpOnly: conf.HTTPOnly,
	}
	store.Options(*strOptions)
}

// SetUserSession sets user session
func SetUserSession(c *gin.Context, userID int) {
	session := sessions.Default(c)
	v := session.Get("uid")
	if v == nil {
		session.Set("uid", userID)
		session.Save()
	}
}

// IsLogin returns boolean whether user have already login or not and uid
func IsLogin(c *gin.Context) (bool, int) {
	session := sessions.Default(c)
	v := session.Get("uid")
	if v == nil {
		return false, 0
	}
	return true, v.(int)
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
		err = errors.New("token is not allowed as blank")
	} else if GetTokenSession(c) == "" {
		err = errors.New("token is missing. Session might have expired")
	} else if GetTokenSession(c) != token {
		err = errors.New("token is invalid")
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
