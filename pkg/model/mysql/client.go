package mysql

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/storage/mysql"
)

func newDBModel(env string, conf *config.MySQLContentConfig) (*DBModel, error) {
	ms, err := mysql.NewMySQL(env, conf)
	if err != nil {
		return nil, err
	}

	return &DBModel{
		DB: ms,
	}, nil
}
