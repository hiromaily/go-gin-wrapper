package models

import (
	"time"
)

const msyqlDatetimeFormat = "2006-01-02 15:04:05"

//日付がnullを許容するのであれば、*time.Timeを使用する
//日付をperseするときには、注意が必要
// db, err := sql.Open("mysql", "root:@/?parseTime=true")
// http://stackoverflow.com/questions/29341590/go-parse-time-from-database
type Users struct {
	UserId    uint16    `column:"user_id"          db:"user_id"`
	FirstName string    `column:"first_name"       db:"first_name"`
	LastName  string    `column:"last_name"        db:"last_name"`
	DeleteFlg string    `column:"delete_flg"       db:"delete_flg"`
	Created   time.Time `column:"create_datetime"  db:"create_datetime"`
	Updated   time.Time `column:"update_datetime"  db:"update_datetime"`
}

//time.Now().Format(msyqlDatetimeFormat)
func NewUser(id uint16, firstN string, lastN string) Users {
	return Users{
		UserId:    id,
		FirstName: firstN,
		LastName:  lastN,
		DeleteFlg: "0",
		//Created:   time.Now().Format(msyqlDatetimeFormat),
		//Updated:   time.Now().Format(msyqlDatetimeFormat),
		Created: time.Now(),
		Updated: time.Now(),
	}
}
