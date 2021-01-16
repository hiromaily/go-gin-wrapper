package main

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/logger"
	"github.com/hiromaily/go-gin-wrapper/pkg/reverseproxy"
)

// Registry interface
type Registry interface {
	NewProxyServer() reverseproxy.Server
}

type registry struct {
	logger *zap.Logger
	conf   *config.Config
}

// NewRegistry returns registry interface
func NewRegistry(conf *config.Config) Registry {
	return &registry{
		conf: conf,
	}
}

// NewProxyServer returns Server interface
func (r *registry) NewProxyServer() reverseproxy.Server {
	return reverseproxy.NewServer(
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
