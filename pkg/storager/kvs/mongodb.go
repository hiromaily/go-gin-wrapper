package kvs

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mongo"
	"github.com/hiromaily/golibs/db/mongodb"
	hrk "github.com/hiromaily/golibs/heroku"
)

// MongoRepo is Mongo object
type MongoRepo struct {
	DB *mongodb.MongoInfo
}

func newMongoStorager(conf *configs.Config) (*MongoRepo, error) {
	c := conf.Mongo

	if conf.Environment == "heroku" {
		host, dbname, user, pass, port, err := hrk.GetMongoInfo("")
		if err == nil {
			c.Host = host
			c.DbName = dbname
			c.User = user
			c.Pass = pass
			c.Port = uint16(port)
		}
	}

	mi, err := mongodb.NewInstance(c.Host, c.DbName, c.User, c.Pass, c.Port)
	if err != nil {
		return nil, err
	}
	if c.DbName != "" {
		mi.GetDB(c.DbName)
	}

	return &MongoRepo{
		DB: mi,
	}, nil
}

// CreateDBModel is to crate db model
func (m *MongoRepo) CreateDBModel() mongomodel.MongoModeler {
	return &mongomodel.MongoModel{
		DB: m.DB,
	}
}
