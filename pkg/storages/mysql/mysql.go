package mysql

import (
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	"github.com/hiromaily/golibs/db/mysql"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
)

func NewMySQL(conf *configs.Config) (*mysql.MS, error) {
	ms := &mysql.MS{}

	if conf.Environment == "heroku" {
		//Heroku
		lg.Debug("HEROKU mode")

		//database
		host, dbname, user, pass, err := hrk.GetMySQLInfo("")
		if err != nil {
			return nil, errors.Wrap(err, "fail to call heroku.GetMySQLInfo()")
		}
		ms, err = mysql.NewInstance(host, dbname, user, pass, 3306)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call mysql.NewInstance()")
		}
	} else {
		var err error

		//TODO:For test
		//dbInfo := conf.MySQL.Test

		dbInfo := conf.MySQL
		ms, err = mysql.NewInstance(dbInfo.Host, dbInfo.DbName, dbInfo.User, dbInfo.Pass, dbInfo.Port)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call mysql.NewInstance()")
		}
	}
	ms.DB.SetMaxIdleConns(50)

	return ms, nil
}
