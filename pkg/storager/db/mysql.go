package db

import (
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	dbmodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mysql"
	"github.com/hiromaily/golibs/db/mysql"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
)

// MySQLRepo is MySQL object
type MySQLRepo struct {
	DB *mysql.MS
}

func newDBStorager(conf *configs.Config) (*MySQLRepo, error) {
	//TODO: MySQL or dummy??
	//TODO: how to handle test mode

	ms := &MySQLRepo{}

	if conf.Environment == "heroku" {
		//Heroku
		lg.Debug("HEROKU mode")

		//database
		host, dbname, user, pass, err := hrk.GetMySQLInfo("")
		if err != nil {
			return nil, errors.Wrap(err, "fail to call heroku.GetMySQLInfo()")
		}
		ms.DB, err = mysql.NewInstance(host, dbname, user, pass, 3306)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call mysql.NewInstance()")
		}
	} else {
		var err error

		//TODO:For test
		//dbInfo := conf.MySQL.Test

		dbInfo := conf.MySQL
		ms.DB, err = mysql.NewInstance(dbInfo.Host, dbInfo.DbName, dbInfo.User, dbInfo.Pass, dbInfo.Port)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call mysql.NewInstance()")
		}
	}
	ms.DB.SetMaxIdleConns(50)

	return ms, nil
}

// CreateDBModel is to create db model
func (m *MySQLRepo) CreateDBModel() dbmodel.DBModeler {
	return &dbmodel.DBModel{
		DB: m.DB,
	}
}
