package main

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mongo"
	dbmodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mysql"
	"github.com/hiromaily/go-gin-wrapper/pkg/server"
	"github.com/hiromaily/go-gin-wrapper/pkg/storager/session"
	"github.com/hiromaily/golibs/auth/jwt"
)

// Registry is for registry interface
type Registry interface {
	NewServerer(port int) server.Serverer
}

type registry struct {
	conf *configs.Config
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *configs.Config) Registry {
	return &registry{conf: conf}
}

// NewBooker is to register for booker interface
func (r *registry) NewServerer(port int) server.Serverer {
	r.initAuth()

	return server.NewServerer(
		r.conf,
		port,
		r.newDBModeler(),
		r.newMongoModeler(),
		r.newSessionStorager(),
	)
}

func (r *registry) newDBModeler() dbmodel.DBModeler {
	storager, err := dbmodel.NewDBModeler(r.conf)
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

func (r *registry) newSessionStorager() session.SessionStorager {
	storager, err := session.NewSessionStorager(r.conf)
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
