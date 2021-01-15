package controller

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/model/mongo"
	"github.com/hiromaily/go-gin-wrapper/pkg/repository"
)

// TODO: define interface

// Controller is controller object
type Controller struct {
	logger    *zap.Logger
	userRepo  repository.UserRepositorier
	mongo     mongomodel.MongoModeler
	apiHeader *config.HeaderConfig
	auth      *config.AuthConfig
	cors      *config.CORSConfig
	// TODO: session should be added here
}

// NewController is to return Controller
func NewController(
	userRepo repository.UserRepositorier,
	logger *zap.Logger,
	mongo mongomodel.MongoModeler,
	apiHeader *config.HeaderConfig,
	auth *config.AuthConfig,
	cors *config.CORSConfig) *Controller {
	return &Controller{
		userRepo:  userRepo,
		logger:    logger,
		mongo:     mongo,
		apiHeader: apiHeader,
		auth:      auth,
		cors:      cors,
	}
}
