package controllers

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/model/mongo"
	dbmodel "github.com/hiromaily/go-gin-wrapper/pkg/model/mysql"
)

// TODO: define interface

// Controller is controller object
type Controller struct {
	db        dbmodel.DBModeler
	mongo     mongomodel.MongoModeler
	apiHeader *config.HeaderConfig
	auth      *config.AuthConfig
	cors      *config.CORSConfig
	// TODO: session should be added here
}

// NewController is to return Controller
func NewController(
	db dbmodel.DBModeler,
	mongo mongomodel.MongoModeler,
	apiHeader *config.HeaderConfig,
	auth *config.AuthConfig,
	cors *config.CORSConfig) *Controller {
	return &Controller{
		db:        db,
		mongo:     mongo,
		apiHeader: apiHeader,
		auth:      auth,
		cors:      cors,
	}
}
