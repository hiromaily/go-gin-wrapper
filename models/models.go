package models

import (
	"github.com/hiromaily/golibs/db/mysql"
)

//http://qiita.com/tenntenn/items/e04441a40aeb9c31dbaf
//http://blog.monochromegane.com/blog/2014/03/23/struct-implementaion-patterns-in-golang/

const msyqlDatetimeFormat = "2006-01-02 15:04:05"

//extention of mysql.DBInfo
type Models struct {
	Db *mysql.MS
}

var db Models

func New() {
	//db = &Models{}
	db = Models{}
	db.Db = mysql.GetDBInstance()
}

//using singleton design pattern
func GetDBInstance() *Models {
	if db.Db == nil {
		New()
	}
	return &db
}
