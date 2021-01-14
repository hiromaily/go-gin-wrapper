package mongo

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/storage/mongodb"
)

func newMongoModel(conf *config.Config) (*MongoModel, error) {
	mi, err := mongodb.NewMongo(conf)
	if err != nil {
		return nil, err
	}

	return &MongoModel{
		DB: mi,
	}, nil
}
