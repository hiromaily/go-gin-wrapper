package controllers

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mongo"
	dbmodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mysql"
)

//TODO: define interface

// Controller is controller object
type Controller struct {
	db        dbmodel.DBModeler
	mongo     mongomodel.MongoModeler
	apiHeader *configs.HeaderConfig
	auth      *configs.AuthConfig
	cors      *configs.CORSConfig
	//TODO: session should be added here
}

// NewController is to return Controller
func NewController(
	db dbmodel.DBModeler,
	mongo mongomodel.MongoModeler,
	apiHeader *configs.HeaderConfig,
	auth *configs.AuthConfig,
	cors *configs.CORSConfig) *Controller {

	return &Controller{
		db:        db,
		mongo:     mongo,
		apiHeader: apiHeader,
		auth:      auth,
		cors:      cors,
	}
}
