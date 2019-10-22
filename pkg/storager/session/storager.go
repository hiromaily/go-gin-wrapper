package session

import "github.com/hiromaily/go-gin-wrapper/pkg/configs"

// SessionStorager is Storager interface
type SessionStorager interface {
	SetSession(string) (bool, error)
}

// NewSessionStorager is to return KVSStorager interface
func NewSessionStorager(conf *configs.Config) (SessionStorager, error) {
	//logic is here, if switch is required
	return newRedisStorager(conf)
}
