package redis

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"

	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
)

func NewRedis(conf *configs.Config) (*redis.Conn, error) {
	red := conf.Redis

	var conn redis.Conn
	var err error
	if conf.Environment == "heroku" {
		conn, err = redis.DialURL(os.Getenv("REDIS_URL"))
	} else {
		if red.Pass != "" {
			//plus password
			conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", red.Host, red.Port), redis.DialPassword(red.Pass))
		} else {
			conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", red.Host, red.Port))
		}
	}
	if err != nil {
		return nil, err
	}

	return &conn, nil
}
