package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
)

// NewMySQL creates mysql db connection
func NewMySQL(conf *config.MySQLContentConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4",
			conf.User,
			conf.Pass,
			conf.Host,
			conf.DBName))
	if err != nil {
		return nil, errors.Errorf("Connection(): error: %v", err)
	}
	return db, nil
}
