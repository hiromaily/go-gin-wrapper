package ginsession

import (
	"github.com/gin-gonic/contrib/sessions"
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
