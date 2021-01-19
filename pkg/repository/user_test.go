package repository

import (
	"database/sql"
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/storage/mysql"
	"testing"
)

func setup(){

}

func newMySQLClient(dbConf *config.MySQLContent) *sql.DB {
	dbConn, err := mysql.NewMySQL(dbConf)
	if err != nil {
		panic(err)
	}
	return dbConn
}

func TestIsUserEmail(t *testing.T) {


}
