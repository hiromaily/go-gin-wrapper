package controller

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/repository"
)

// Controller  interface
type Controller interface {
	Acounter
	Adminer
	APIJWTer
	APIUser
	APILister
	Baser
	Chater
	Errorer
	Loginer
	OAuther
}

// controller is controller object
type controller struct {
	logger    *zap.Logger
	userRepo  repository.UserRepositorier
	apiHeader *config.Header
	cors      *config.CORS
	auth      *config.Auth
	// TODO: session should be added here
}

// NewController is to return Controller
func NewController(
	logger *zap.Logger,
	userRepo repository.UserRepositorier,
	apiHeader *config.Header,
	cors *config.CORS,
	auth *config.Auth) Controller {
	return &controller{
		logger:    logger,
		userRepo:  userRepo,
		apiHeader: apiHeader,
		auth:      auth,
		cors:      cors,
	}
}
