package mysql

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	"github.com/hiromaily/go-gin-wrapper/pkg/storages/mysql"
)

func newDBModel(conf *configs.Config) (*DBModel, error) {
	ms, err := mysql.NewMySQL(conf)
	if err != nil {
		return nil, err
	}

	return &DBModel{
		DB: ms,
	}, nil
}
