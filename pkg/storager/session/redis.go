package session

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"

	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
)

//TODO: integrate pkg/libs/ginsession/

// RedisRepo is Redis object
type RedisRepo struct {
	RD *redis.Conn
}

func newRedisStorager(conf *configs.Config) (*RedisRepo, error) {
	red := conf.Redis
	//if os.Getenv("HEROKU_FLG") == "1" {
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
		//TODO:
		//red.Session = true
	}
	if err != nil {
		return nil, err
	}

	return &RedisRepo{
		RD: &conn,
	}, nil

}

// SetSession is to set session
func (r *RedisRepo) SetSession(string) (bool, error) {
	//if red.Session && red.Host != "" && red.Port != 0 {
	//	lg.Debug("redis session start")
	//	sess.SetSession(r, fmt.Sprintf("%s:%d", red.Host, red.Port), red.Pass)
	//} else {
	//	sess.SetSession(r, "", "")
	//}
	return false, nil
}
