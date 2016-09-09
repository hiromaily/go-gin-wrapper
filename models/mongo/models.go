package mongo

import (
	mongo "github.com/hiromaily/golibs/db/mongodb"
)

//
//extension of db/mysql package.
//

// Models is extension of mongo.MongoInfo
type Models struct {
	Db *mongo.MongoInfo
}

var db Models

// when making mongo instance, first you should use mongo.New()
func new() {
	db = Models{}
	db.Db = mongo.GetMongo()
}

// GetDB is to get mongo instance. it's using singleton design pattern.
func GetDB() *Models {
	if db.Db == nil {
		new()
	}
	return &db
}
