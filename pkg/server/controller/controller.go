package controller

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/auth/jwts"
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
	userRepo  repository.UserRepository
	jwter     jwts.JWTer
	apiHeader *config.Header
	cors      *config.CORS
	auth      *config.Auth
}

// NewController is to return Controller
func NewController(
	logger *zap.Logger,
	userRepo repository.UserRepository,
	jwter jwts.JWTer,
	apiHeader *config.Header,
	cors *config.CORS,
	auth *config.Auth) Controller {
	return &controller{
		logger:    logger,
		userRepo:  userRepo,
		jwter:     jwter,
		apiHeader: apiHeader,
		auth:      auth,
		cors:      cors,
	}
}
