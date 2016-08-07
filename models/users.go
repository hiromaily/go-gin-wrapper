package models

import (
	"fmt"
	hs "github.com/hiromaily/golibs/cipher/hash"
	lg "github.com/hiromaily/golibs/log"
	u "github.com/hiromaily/golibs/utils"
)

//If use nil on Date column, set *time.Time
//Be careful when parsing is required on Date type
// e.g. db, err := sql.Open("mysql", "root:@/?parseTime=true")
// http://stackoverflow.com/questions/29341590/go-parse-time-from-database
type Users struct {
	//UserId    int       `column:"user_id"          db:"user_id"`
	FirstName string `column:"first_name"       db:"first_name"`
	LastName  string `column:"last_name"        db:"last_name"`
	Email     string `column:"email"            db:"email"`
	Password  string `column:"password"         db:"password"`
	//DeleteFlg string    `column:"delete_flg"       db:"delete_flg"`
	//Created   time.Time `column:"create_datetime"  db:"create_datetime"`
	//Updated   time.Time `column:"update_datetime"  db:"update_datetime"`
}

//User authorization when login
func (us *Models) IsUserEmail(email string, password string) (int, error) {
	sql := "SELECT user_id, email, password FROM t_users WHERE email=? AND delete_flg=?"
	data, _, err := us.Db.SelectSQLAllField(sql, email, 0)
	lg.Debugf("mysql data %+v", data)

	if err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, fmt.Errorf("email may be wrong.")
	}

	if u.Itos(data[0]["password"]) != hs.GetMD5Plus(password, "") {
		lg.Debugf("data[0]['password'] : %s\n", u.Itos(data[0]["password"]))
		lg.Debugf("data[0]['email'] : %s\n", u.Itos(data[0]["email"]))
		lg.Debugf("GetMD5Plus : %s\n", hs.GetMD5Plus(password, u.Itos(data[0]["email"])))

		return 0, fmt.Errorf("password is invalid.")
	}
	//
	return u.Itoi(data[0]["user_id"]), nil
}

// Get User List
func (us *Models) GetUserList(id string) (data []map[string]interface{}, err error) {
	sql := "SELECT user_id, first_name, last_name FROM t_users WHERE delete_flg=?"
	if id != "" {
		sql += " AND user_id=?"
		data, _, err = us.Db.SelectSQLAllField(sql, 0, u.Atoi(id))
	} else {
		//data, _, err := us.Db.SelectSQLAllField(sql, 1)
		data, _, err = us.Db.SelectSQLAllField(sql, 0)
	}
	if err != nil {
		lg.Errorf("SQL may be wrong. : %s\n", err.Error())
		return nil, err
	} else if len(data) == 0 {
		lg.Info("No data.")
		return nil, nil
	}
	return data, nil
}

// Insert User
func (us *Models) InsertUser(users *Users) (int64, error) {
	sql := "INSERT INTO t_users (first_name, last_name, email, password) VALUES (?,?,?,?)"

	//hash password
	return us.Db.InesrtSQL(sql, users.FirstName, users.LastName, users.Email, hs.GetMD5Plus(users.Password, ""))
}

// Update User
func (us *Models) UpdateUser(users *Users, id string) (int64, error) {
	vals := []interface{}{}
	sql := "UPDATE t_users SET"
	if users.FirstName != "" {
		sql += " first_name=?,"
		vals = append(vals, users.FirstName)
	}
	if users.LastName != "" {
		sql += " last_name=?,"
		vals = append(vals, users.LastName)
	}
	if users.Email != "" {
		sql += " email=?,"
		vals = append(vals, users.Email)
	}
	if users.Password != "" {
		sql += " password=?,"
		vals = append(vals, hs.GetMD5Plus(users.Password, ""))
	}
	//remove last comma
	sql = string(sql[:(len(sql) - 1)])

	//user id
	sql += " WHERE user_id=?"
	vals = append(vals, u.Atoi(id))

	//sql debug
	//lg.Debugf("update sql: %s", sql)

	return us.Db.UpdateSQL(sql, vals...)
}

// Delete User
func (us *Models) DeleteUser(id string) error {
	sql := "DELETE FROM t_users WHERE user_id=?"
	return us.Db.ExecSQL(sql, u.Atoi(id))
}

/*
//time.Now().Format(msyqlDatetimeFormat)
func NewUser(id int, firstN string, lastN string) Users {
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
*/
