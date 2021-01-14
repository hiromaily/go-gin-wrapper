package mysql

import (
	"fmt"

	str "github.com/hiromaily/go-gin-wrapper/pkg/strings"
	hs "github.com/hiromaily/golibs/cipher/hash"
	"github.com/hiromaily/golibs/db/mysql"
	lg "github.com/hiromaily/golibs/log"
)

// DBModel is extension of mysql.DBInfo
type DBModel struct {
	DB *mysql.MS
}

// IsUserEmail is for User authorization when trying login
func (us *DBModel) IsUserEmail(email string, password string) (int, error) {
	type LoginUser struct {
		ID       int
		Email    string
		Password string
	}

	sql := "SELECT user_id, email, password FROM t_users WHERE email=? AND delete_flg=? LIMIT 1"

	var user LoginUser
	b := us.DB.SelectIns(sql, email, 0).ScanOne(&user)

	if us.DB.Err != nil {
		lg.Debugf("IsUserEmail() Error: %s", us.DB.Err)
		return 0, us.DB.Err
	}

	// no data
	if !b {
		return 0, fmt.Errorf("email may be wrong")
	}

	if user.Password != hs.GetMD5Plus(password, "") {
		return 0, fmt.Errorf("password is invalid")
	}
	return user.ID, nil
}

// OAuth2Login is for OAuth2 login
func (us *DBModel) OAuth2Login(email string) (*UserAuth, error) {
	// 0:no user -> register and login
	// 1:existing user (google) -> login
	// 2:existing user (no auth or another auth) -> err
	sql := "SELECT user_id, oauth2_flg FROM t_users WHERE email=? AND delete_flg=?"

	var user UserAuth

	b := us.DB.SelectIns(sql, email, 0).ScanOne(&user)
	if us.DB.Err != nil {
		// 0:no user -> register and login
		return nil, us.DB.Err
	}

	// no data
	if !b {
		return nil, nil
	}

	return &user, nil
}

// GetUserIds is to get user IDs
func (us *DBModel) GetUserIds(users interface{}) error {
	sql := "SELECT user_id FROM t_users WHERE delete_flg=?"

	us.DB.SelectIns(sql, 0).Scan(users)

	if us.DB.Err != nil {
		return us.DB.Err
	}

	return nil
}

// GetUserList is to get user list
func (us *DBModel) GetUserList(users interface{}, id string) (bool, error) {
	// lg.Debug(mysql.ColumnForSQL(users))

	// remove password
	// fields := strings.Replace(mysql.ColumnForSQL(users), "password,", "", 1)
	// lg.Debug(fields)

	sql := "SELECT %s FROM t_users WHERE delete_flg=?"
	sql = fmt.Sprintf(sql, mysql.ColumnForSQL(users))

	// TODO: When Test for result is 0 record, set 1 to delFlg
	delFlg := 0

	var b bool
	if id != "" {
		sql += " AND user_id=?"
		b = us.DB.SelectIns(sql, delFlg, str.Atoi(id)).ScanOne(users)
	} else {
		b = us.DB.SelectIns(sql, delFlg).Scan(users)
	}

	if us.DB.Err != nil {
		return false, us.DB.Err
	}

	return b, nil
}

// InsertUser is to insert user
func (us *DBModel) InsertUser(users *Users) (int64, error) {
	lg.Debug(mysql.ColumnForSQL(users))

	sql := "INSERT INTO t_users (first_name, last_name, email, password) VALUES (?,?,?,?)"
	// sql = fmt.Sprintf(sql, mysql.ColumnForSQL(users))
	if users.OAuth2Flg != "" {
		sql := "INSERT INTO t_users (first_name, last_name, email, password, oauth2_flg) VALUES (?,?,?,?,?)"
		// hash password
		return us.DB.Insert(sql, users.FirstName, users.LastName, users.Email, hs.GetMD5Plus(users.Password, ""), users.OAuth2Flg)
	}

	// hash password
	return us.DB.Insert(sql, users.FirstName, users.LastName, users.Email, hs.GetMD5Plus(users.Password, ""))
}

// UpdateUser is to update user
func (us *DBModel) UpdateUser(users *Users, id string) (int64, error) {
	// lg.Debug(mysql.ColumnForSQL(users))

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
	if users.Updated != "" {
		sql += " update_datetime=?,"
		vals = append(vals, users.Updated)
	}

	// remove last comma
	sql = string(sql[:(len(sql) - 1)])

	// user id
	sql += " WHERE user_id=?"
	vals = append(vals, str.Atoi(id))

	// sql debug
	// lg.Debugf("update sql: %s", sql)

	// Update
	return us.DB.Exec(sql, vals...)
}

// DeleteUser is to delete user
func (us *DBModel) DeleteUser(id string) (int64, error) {
	sql := "DELETE FROM t_users WHERE user_id=?"
	return us.DB.Exec(sql, str.Atoi(id))
}
