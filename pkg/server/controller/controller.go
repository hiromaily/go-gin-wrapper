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
	cors      *config.CORSConfig
	auth      *config.AuthConfig
	// TODO: session should be added here
}

// NewController is to return Controller
func NewController(
	logger *zap.Logger,
	userRepo repository.UserRepositorier,
	mongo mongomodel.MongoModeler,
	apiHeader *config.HeaderConfig,
	cors *config.CORSConfig,
	auth *config.AuthConfig) *Controller {
	return &Controller{
		logger:    logger,
		userRepo:  userRepo,
		mongo:     mongo,
		apiHeader: apiHeader,
		auth:      auth,
		cors:      cors,
	}
}
