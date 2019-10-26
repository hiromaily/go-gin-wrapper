package mongo

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	"github.com/hiromaily/go-gin-wrapper/pkg/storages/mongodb"
)

func newMongoModel(conf *configs.Config) (*MongoModel, error) {
	mi, err := mongodb.NewMongo(conf)
	if err != nil {
		return nil, err
	}

	return &MongoModel{
		DB: mi,
	}, nil
}
