package mysql

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	"github.com/hiromaily/go-gin-wrapper/pkg/storages/mysql"
)

func newDBModel(env string, conf *configs.MySQLContentConfig) (*DBModel, error) {
	ms, err := mysql.NewMySQL(env, conf)
	if err != nil {
		return nil, err
	}

	return &DBModel{
		DB: ms,
	}, nil
}
