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
	Id        int    `column:"user_id"`
	FirstName string `column:"first_name"`
	LastName  string `column:"last_name"`
	Email     string `column:"email"`
	Password  string `column:"password"`
	//DeleteFlg string    `column:"delete_flg"       db:"delete_flg"`
	//Created   time.Time `column:"create_datetime"  db:"create_datetime"`
	//Updated   time.Time `column:"update_datetime"  db:"update_datetime"`
}

//User authorization when login
func (us *Models) IsUserEmail(email string, password string) (int, error) {
	type LoginUser struct {
		Id       int
		Email    string
		Password string
	}

	sql := "SELECT user_id, email, password FROM t_users WHERE email=? AND delete_flg=? LIMIT 1"

	var user LoginUser
	b := us.Db.SelectIns(sql, email, 0).ScanOne(&user)

	if us.Db.Err != nil {
		lg.Debugf("IsUserEmail() Error: %s", us.Db.Err)
		return 0, us.Db.Err
	}

	//no data
	if !b {
		return 0, fmt.Errorf("email may be wrong.")
	}

	if user.Password != hs.GetMD5Plus(password, "") {
		return 0, fmt.Errorf("password is invalid.")
	}
	return user.Id, nil
}

// Get User List
func (us *Models) GetUserList(users interface{}, id, sql string) (bool, error) {
	defaultSql := "SELECT user_id, first_name, last_name, email, password FROM t_users WHERE delete_flg=?"
	if sql == "" {
		sql = defaultSql
	}

	//TODO: When Test for result is 0 record, set 1 to delFlg
	delFlg := 0

	var b bool
	if id != "" {
		sql += " AND user_id=?"
		b = us.Db.SelectIns(sql, delFlg, u.Atoi(id)).ScanOne(users)
	} else {
		b = us.Db.SelectIns(sql, delFlg).Scan(users)
	}

	if us.Db.Err != nil {
		return false, us.Db.Err
	}

	return b, nil
}

// Insert User
func (us *Models) InsertUser(users *Users) (int64, error) {
	sql := "INSERT INTO t_users (first_name, last_name, email, password) VALUES (?,?,?,?)"

	//hash password
	return us.Db.Insert(sql, users.FirstName, users.LastName, users.Email, hs.GetMD5Plus(users.Password, ""))
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

	//Update
	return us.Db.Exec(sql, vals...)
}

// Delete User
func (us *Models) DeleteUser(id string) (int64, error) {
	sql := "DELETE FROM t_users WHERE user_id=?"
	return us.Db.Exec(sql, u.Atoi(id))
}
