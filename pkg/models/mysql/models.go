package mysql

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
)

// DBModeler is DBModeler interface
type DBModeler interface {
	IsUserEmail(email string, password string) (int, error)
	OAuth2Login(email string) (*UserAuth, error)
	GetUserIds(users interface{}) error
	GetUserList(users interface{}, id string) (bool, error)
	InsertUser(users *Users) (int64, error)
	UpdateUser(users *Users, id string) (int64, error)
	DeleteUser(id string) (int64, error)
}

// NewDBModeler is to return DBModeler interface
func NewDBModeler(conf *configs.Config) (DBModeler, error) {
	//logic is here, if switching is required

	//MongoDB
	return newDBModel(conf)

	//or dummy
	//return &DummyMySQL{}, nil
}
