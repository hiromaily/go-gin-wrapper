package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/auth/jwt"
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/logger"
	"github.com/hiromaily/go-gin-wrapper/pkg/repository"
	"github.com/hiromaily/go-gin-wrapper/pkg/server"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/controller"
	"github.com/hiromaily/go-gin-wrapper/pkg/storage/mysql"
)

// Registry is for registry interface
type Registry interface {
	NewServer() server.Server
}

type registry struct {
	conf        *config.Config
	logger      *zap.Logger
	gin         *gin.Engine
	isTestMode  bool
	mysqlClient *sql.DB
	// redisClient *redis.Conn
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Config, isTestMode bool) Registry {
	return &registry{
		isTestMode: isTestMode,
		conf:       conf,
	}
}

func (r *registry) initAuth() {
	auth := r.conf.API.JWT
	if auth.Mode == jwt.HMAC && auth.Secret != "" {
		jwt.InitSecretKey(auth.Secret)
	} else if auth.Mode == jwt.RSA && auth.PrivateKey != "" && auth.PublicKey != "" {
		err := jwt.InitKeys(auth.PrivateKey, auth.PublicKey)
		if err != nil {
			panic(err)
		}
	} else {
		jwt.InitEncrypted(jwt.HMAC)
	}
}

// NewServerer is to register for serverer interface
func (r *registry) NewServer() server.Server {
	r.initAuth()

	return server.NewServer(
		r.newGin(),
		r.newMiddleware(),
		r.newController(),
		r.newLogger(),
		r.newUserRepository(),
		r.conf,
		r.isTestMode,
	)
}

func (r *registry) newGin() *gin.Engine {
	if r.gin == nil {
		r.gin = gin.New()
	}
	return r.gin
}

func (r *registry) newMiddleware() server.Middlewarer {
	return server.NewMiddleware(
		r.newLogger(),
		r.conf.Server,
		r.conf.Proxy,
		r.conf.API,
		r.conf.API.CORS,
		r.conf.Develop,
	)
}

func (r *registry) newController() *controller.Controller {
	return controller.NewController(
		r.newLogger(),
		r.newUserRepository(),
		r.conf.API.Header,
		r.conf.API.CORS,
		r.conf.Auth,
	)
}

func (r *registry) newLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(r.conf.Logger)
	}
	return r.logger
}

func (r *registry) newMySQLClient() *sql.DB {
	if r.mysqlClient == nil {
		dbConf := r.conf.MySQL.MySQLContentConfig
		if r.isTestMode {
			dbConf = r.conf.MySQL.Test
		}

		dbConn, err := mysql.NewMySQL(dbConf)
		if err != nil {
			panic(err)
		}
		r.mysqlClient = dbConn
	}
	if r.conf.MySQL.IsDebugLog {
		boil.DebugMode = true
	}

	return r.mysqlClient
}

func (r *registry) newUserRepository() repository.UserRepositorier {
	return repository.NewUserRepository(r.newMySQLClient(), r.newLogger())
}

// newRedisConn is to create redis connection
//func (r *registry) newRedisConn(conf *config.RedisConfig) *redis.Conn {
//	if r.redisClient == nil {
//		redisConn, err := rd.NewRedis(conf)
//		if err != nil {
//			panic(err)
//		}
//		r.redisClient = redisConn
//	}
//	return r.redisClient
//}
