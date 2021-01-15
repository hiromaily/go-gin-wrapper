package main

import (
	"database/sql"
	"log"

	"github.com/garyburd/redigo/redis"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/auth/jwt"
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/logger"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/model/mongo"
	"github.com/hiromaily/go-gin-wrapper/pkg/repository"
	"github.com/hiromaily/go-gin-wrapper/pkg/server"
	"github.com/hiromaily/go-gin-wrapper/pkg/storage/mysql"
	rd "github.com/hiromaily/go-gin-wrapper/pkg/storage/redis"
)

// Registry is for registry interface
type Registry interface {
	NewServerer(port int) server.Serverer
}

type registry struct {
	isTestMode  bool
	conf        *config.Config
	logger      *zap.Logger
	mysqlClient *sql.DB
	redisConn   *redis.Conn
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Config, isTestMode bool) Registry {
	return &registry{
		isTestMode: isTestMode,
		conf:       conf,
		redisConn:  newRedisConn(conf),
	}
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
	//if r.conf.MySQL.Debug {
	//	boil.DebugMode = true
	//}

	return r.mysqlClient
}

func (r *registry) newUserRepository() repository.UserRepositorier {
	return repository.NewUserRepository(r.newMySQLClient(), r.newLogger())
}

// newRedisConn is to create redis connection
func newRedisConn(conf *config.Config) *redis.Conn {
	conn, err := rd.NewRedis(conf)
	if err != nil {
		log.Println("failed to create redis connection")
		return nil
	}
	return conn
}

// NewServerer is to register for serverer interface
func (r *registry) NewServerer(port int) server.Serverer {
	r.initAuth()

	return server.NewServerer(
		r.isTestMode,
		r.conf,
		port,
		r.newUserRepository(),
		r.newMongoModeler(),
	)
}

func (r *registry) newMongoModeler() mongomodel.MongoModeler {
	storager, err := mongomodel.NewMongoModeler(r.conf)
	if err != nil {
		panic(err)
	}
	return storager
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
		// lg.Debug("JWT Auth is not available because of toml settings.")
	}
}
