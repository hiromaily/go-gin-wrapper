package main

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	"github.com/hiromaily/go-gin-wrapper/pkg/reverseproxy"
)

// Registry is for registry interface
type Registry interface {
	NewProxyServerer() reverseproxy.Serverer
}

type registry struct {
	conf *configs.Config
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *configs.Config) Registry {
	return &registry{
		conf: conf,
	}
}

// NewProxyServerer is to register for serverer interface
func (r *registry) NewProxyServerer() reverseproxy.Serverer {
	return reverseproxy.NewServerer(
		r.conf,
	)
}
