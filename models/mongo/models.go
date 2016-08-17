package mongo

import (
	mongo "github.com/hiromaily/golibs/db/mongodb"
)

//
//extension of db/mysql package.
//

//extention of mysql.DBInfo
type Models struct {
	Db *mongo.MongoInfo
}

var db Models

func new() {
	db = Models{}
	db.Db = mongo.GetMongo()
}

//using singleton design pattern
func GetDB() *Models {
	if db.Db == nil {
		new()
	}
	return &db
}
