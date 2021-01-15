package main

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/logger"
	"github.com/hiromaily/go-gin-wrapper/pkg/reverseproxy"
)

// Registry is for registry interface
type Registry interface {
	NewProxyServerer() reverseproxy.Serverer
}

type registry struct {
	logger *zap.Logger
	conf   *config.Config
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Config) Registry {
	return &registry{
		conf: conf,
	}
}

// NewProxyServerer is to register for serverer interface
func (r *registry) NewProxyServerer() reverseproxy.Serverer {
	return reverseproxy.NewServerer(
		r.newLogger(),
		r.conf,
	)
}

func (r *registry) newLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(r.conf.Logger)
	}
	return r.logger
}
