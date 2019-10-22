package kvs

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mongo"
)

// KVSStorager is KVSStorager interface
type KVSStorager interface {
	CreateDBModel() mongomodel.MongoModeler
}

// NewKVSStorager is to return KVSStorager interface
func NewKVSStorager(conf *configs.Config) (KVSStorager, error) {
	//logic is here, if switch is required

	//MongoDB
	return newMongoStorager(conf)
}
