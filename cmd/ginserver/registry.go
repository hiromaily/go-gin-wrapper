package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/auth/jwts"
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
	"github.com/hiromaily/go-gin-wrapper/pkg/heroku"
	"github.com/hiromaily/go-gin-wrapper/pkg/logger"
	"github.com/hiromaily/go-gin-wrapper/pkg/repository"
	"github.com/hiromaily/go-gin-wrapper/pkg/server"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/controller"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader/cors"
	"github.com/hiromaily/go-gin-wrapper/pkg/storage/mysql"
	"github.com/hiromaily/go-gin-wrapper/pkg/token"
)

// Registry interface
type Registry interface {
	NewServer() server.Server
}

type registry struct {
	conf        *config.Root
	logger      *zap.Logger
	gin         *gin.Engine
	session     ginsession.Sessioner
	jwter       jwts.JWTer
	token       token.Generator
	isTestMode  bool
	mysqlClient *sql.DB
	hash        encryption.MD5
	// redisClient *redis.Conn
}

// NewRegistry returns registry interface
func NewRegistry(conf *config.Root, isTestMode bool) Registry {
	return &registry{
		isTestMode: isTestMode,
		conf:       conf,
	}
}

// NewServer returns Server interface
func (r *registry) NewServer() server.Server {
	return server.NewServer(
		r.newGin(),
		r.newSessionStore(),
		r.newMiddleware(),
		r.newController(),
		r.newLogger(),
		r.newMySQLClient(),
		r.newUserRepository(),
		r.conf,
		r.isTestMode,
	)
}

func (r *registry) newGin() *gin.Engine {
	if r.gin == nil {
		r.gin = gin.New()
		if r.conf.Server.IsRelease {
			gin.SetMode(gin.ReleaseMode)
		}
	}
	return r.gin
}

func (r *registry) newSessionStore() sessions.Store {
	red := r.conf.Redis
	herokuRedisSetting(red)

	if red.IsSession && red.Host != "" && red.Port != 0 {
		r.newLogger().Debug("newSessionStore(): redis session start")
		return ginsession.NewRedisStore(
			r.newLogger(),
			fmt.Sprintf("%s:%d", red.Host, red.Port),
			red.Pass,
			r.conf.Server.Session)
	}
	return ginsession.NewCookieStore(r.conf.Server.Session)
}

func herokuRedisSetting(redisConf *config.Redis) {
	if redisConf.IsHeroku {
		host, pass, port, err := heroku.GetRedisInfo("")
		if err == nil && host != "" && port != 0 {
			redisConf.IsSession = true
			redisConf.Host = host
			redisConf.Port = uint16(port)
			redisConf.Pass = pass
		}
	}
}

func (r *registry) newMiddleware() server.Middlewarer {
	return server.NewMiddleware(
		r.newLogger(),
		r.newSession(),
		r.newJWT(),
		r.newCORS(),
		r.rejectIPs(),
		r.conf.Server,
		r.conf.Proxy,
		r.conf.API,
		r.conf.Develop,
	)
}

func (r *registry) newCORS() cors.CORSer {
	return cors.NewCORS(r.newLogger(), r.conf.API.CORS)
}

// TODO: add logic
func (r *registry) rejectIPs() []string {
	return []string{}
}

func (r *registry) newController() controller.Controller {
	return controller.NewController(
		r.newLogger(),
		r.newUserRepository(),
		r.newSession(),
		r.newJWT(),
		r.conf.API.Header,
		r.conf.Auth,
	)
}

func (r *registry) newSession() ginsession.Sessioner {
	if r.session == nil {
		r.session = ginsession.NewSessioner(r.newLogger(), r.newTokenGenerator())
	}
	return r.session
}

func (r *registry) newJWT() jwts.JWTer {
	if r.jwter == nil {
		auth := r.conf.API.JWT
		var signAlgo jwts.SigAlgoer
		if auth.Mode == jwts.HMAC && auth.Secret != "" {
			signAlgo = jwts.NewHMAC(auth.Secret)
		} else if auth.Mode == jwts.RSA && auth.PrivateKey != "" && auth.PublicKey != "" {
			var err error
			signAlgo, err = jwts.NewRSA(auth.PrivateKey, auth.PublicKey)
			if err != nil {
				panic(err)
			}
		} else {
			panic(errors.New("invalid jwt config"))
		}
		r.jwter = jwts.NewJWT(auth.Audience, signAlgo)
	}
	return r.jwter
}

func (r *registry) newTokenGenerator() token.Generator {
	if r.token == nil {
		r.token = token.NewGenerator(r.conf.Server.Token.Salt)
	}
	return r.token
}

func (r *registry) newLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(r.conf.Logger)
	}
	return r.logger
}

func (r *registry) newMySQLClient() *sql.DB {
	if r.mysqlClient == nil {
		dbConf := r.conf.MySQL.MySQLContent
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

func (r *registry) newUserRepository() repository.UserRepository {
	return repository.NewUserRepository(
		r.newMySQLClient(),
		r.newLogger(),
		r.newHash(),
	)
}

func (r *registry) newHash() encryption.Hasher {
	if r.hash == nil {
		r.hash = encryption.NewMD5(r.conf.Hash.Salt1, r.conf.Hash.Salt2)
	}
	return r.hash
}

// newRedisConn returns redis connection
//func (r *registry) newRedisConn(conf *config.Redis) *redis.Conn {
//	if r.redisClient == nil {
//		redisConn, err := rds.NewRedis(conf)
//		if err != nil {
//			panic(err)
//		}
//		r.redisClient = redisConn
//	}
//	return r.redisClient
//}
