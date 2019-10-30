package mysql

import (
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	"github.com/hiromaily/golibs/db/mysql"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
)

// NewMySQL is to return mysql connection
func NewMySQL(env string, conf *configs.MySQLContentConfig) (*mysql.MS, error) {
	var ms *mysql.MS

	if env == "heroku" {
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

		ms, err = mysql.NewInstance(conf.Host, conf.DbName, conf.User, conf.Pass, conf.Port)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call mysql.NewInstance()")
		}
	}
	ms.DB.SetMaxIdleConns(50)

	return ms, nil
}
