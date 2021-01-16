package redis

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
)

// NewRedis is to return redis connection
func NewRedis(conf *config.Redis) (*redis.Conn, error) {
	var conn redis.Conn
	var err error
	if conf.IsHeroku {
		conn, err = redis.DialURL(os.Getenv("REDIS_URL"))
	} else {
		if conf.Pass != "" {
			// plus password
			conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port), redis.DialPassword(conf.Pass))
		} else {
			conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
		}
	}
	if err != nil {
		return nil, err
	}

	return &conn, nil
}
