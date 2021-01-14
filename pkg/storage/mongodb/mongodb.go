package mongodb

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	mdb "github.com/hiromaily/golibs/db/mongodb"
	hrk "github.com/hiromaily/golibs/heroku"
)

// NewMongo is to return mongodb connection
func NewMongo(conf *config.Config) (*mdb.MongoInfo, error) {
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

	mi, err := mdb.NewInstance(c.Host, c.DbName, c.User, c.Pass, c.Port)
	if err != nil {
		return nil, err
	}
	if c.DbName != "" {
		mi.GetDB(c.DbName)
	}

	return mi, err
}
