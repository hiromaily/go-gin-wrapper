package mongodb

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	mdb "github.com/hiromaily/golibs/db/mongodb"
)

// NewMongo is to return mongodb connection
func NewMongo(conf *config.Config) (*mdb.MongoInfo, error) {
	c := conf.Mongo

	mi, err := mdb.NewInstance(c.Host, c.DbName, c.User, c.Pass, c.Port)
	if err != nil {
		return nil, err
	}
	if c.DbName != "" {
		mi.GetDB(c.DbName)
	}

	return mi, err
}
