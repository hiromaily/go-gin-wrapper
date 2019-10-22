package controllers

import (
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mongo"
	dbmodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mysql"
)

//TODO: define interface

// Controller is controller object
type Controller struct {
	db    dbmodel.DBModeler
	mongo mongomodel.MongoModeler
	//TODO: session should be added here
}

// NewController is to return Controller
func NewController(db dbmodel.DBModeler, mongo mongomodel.MongoModeler) *Controller {
	return &Controller{
		db:    db,
		mongo: mongo,
	}
}
