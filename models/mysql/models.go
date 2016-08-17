package mysql

import (
	"github.com/hiromaily/golibs/db/mysql"
)

//
//extension of db/mysql package.
//

//extention of mysql.DBInfo
type Models struct {
	Db *mysql.MS
}

var db Models

func new() {
	db = Models{}
	db.Db = mysql.GetDBInstance()
}

//using singleton design pattern
func GetDB() *Models {
	if db.Db == nil {
		new()
	}
	return &db
}
