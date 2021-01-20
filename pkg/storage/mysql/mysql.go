package mysql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
)

// NewMySQL creates mysql db connection
func NewMySQL(conf *config.MySQLContent) (*sql.DB, error) {
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		conf.User,
		conf.Pass,
		conf.Host,
		conf.Port,
		conf.DBName)
	log.Printf("db source: %s", dbSource)
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		return nil, errors.Errorf("fail to call sql.Open(): %v", err)
	}

	return db, db.Ping()
}
