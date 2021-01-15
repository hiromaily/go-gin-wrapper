package controller

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/repository"
)

// TODO: define interface

// Controller is controller object
type Controller struct {
	logger    *zap.Logger
	userRepo  repository.UserRepositorier
	apiHeader *config.HeaderConfig
	cors      *config.CORSConfig
	auth      *config.AuthConfig
	// TODO: session should be added here
}

// NewController is to return Controller
func NewController(
	logger *zap.Logger,
	userRepo repository.UserRepositorier,
	apiHeader *config.HeaderConfig,
	cors *config.CORSConfig,
	auth *config.AuthConfig) *Controller {
	return &Controller{
		logger:    logger,
		userRepo:  userRepo,
		apiHeader: apiHeader,
		auth:      auth,
		cors:      cors,
	}
}
