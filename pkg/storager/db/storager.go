package db

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	dbmodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mysql"
)

// DBStorager is DBStorager interface
type DBStorager interface {
	CreateDBModel() dbmodel.DBModeler
}

// NewDBStorager is to return DBStorager interface
func NewDBStorager(conf *configs.Config) (DBStorager, error) {
	//logic is here, if switch is required

	//mysql
	return newDBStorager(conf)
}
