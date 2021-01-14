package mongo

import "github.com/hiromaily/go-gin-wrapper/pkg/configs"

// MongoModeler is DBModeler interface
type MongoModeler interface {
	GetNewsData() ([]News, error)
	GetArticlesData(newsID int) ([]Articles, error)
	GetArticlesData2(newsID int) ([]Item2, error)
}

// NewMongoModeler is to return KVSStorager interface
func NewMongoModeler(conf *configs.Config) (MongoModeler, error) {
	// logic is here, if switching is required
	// MongoDB
	return newMongoModel(conf)
	// or dummy
	// return &DummyMongo{}, nil
}
