package main

import (
	"log"

	"github.com/garyburd/redigo/redis"

	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mongo"
	dbmodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mysql"
	"github.com/hiromaily/go-gin-wrapper/pkg/server"
	rd "github.com/hiromaily/go-gin-wrapper/pkg/storages/redis"

	//"github.com/hiromaily/go-gin-wrapper/pkg/storages/redis"
	"github.com/hiromaily/golibs/auth/jwt"
)

// Registry is for registry interface
type Registry interface {
	NewServerer(port int) server.Serverer
}

type registry struct {
	isTestMode bool
	conf       *configs.Config
	redisConn  *redis.Conn
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *configs.Config, isTestMode bool) Registry {
	return &registry{
		isTestMode: isTestMode,
		conf:       conf,
		redisConn:  newRedisConn(conf),
	}
}

// newRedisConn is to create redis connection
func newRedisConn(conf *configs.Config) *redis.Conn {
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
		r.newDBModeler(),
		r.newMongoModeler(),
	)
}

func (r *registry) newDBModeler() dbmodel.DBModeler {
	var dbConf *configs.MySQLContentConfig
	if r.isTestMode {
		dbConf = r.conf.MySQL.Test
	} else {
		dbConf = r.conf.MySQL.MySQLContentConfig
	}

	storager, err := dbmodel.NewDBModeler(r.conf.Environment, dbConf)
	if err != nil {
		panic(err)
	}
	return storager
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
		//lg.Debug("JWT Auth is not available because of toml settings.")
	}
}
