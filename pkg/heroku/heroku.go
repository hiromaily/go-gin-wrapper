package heroku

import (
	"os"
	"strings"

	"github.com/pkg/errors"

	str "github.com/hiromaily/go-gin-wrapper/pkg/strings"
)

// GetRedisInfo returns Redis info from environment variable `REDIS_URL`
func GetRedisInfo(url string) (string, string, int, error) {
	if url == "" {
		if url = os.Getenv("REDIS_URL"); url == "" {
			return "", "", 0, errors.New("`REDIS_URL` is not found")
		}
	}
	//redis://h:xxx@xxx:20819
	//<password>@<hostname>:<port>
	tmp := strings.Split(url, "//")
	if len(tmp) < 2 {
		return "", "", 0, errors.New("url is invalid")
	}
	tmp = strings.Split(tmp[1], ":")
	if len(tmp) < 3 {
		return "", "", 0, errors.New("url is invalid")
	}
	port := str.Atoi(tmp[2])

	tmp = strings.Split(tmp[1], "@")
	if len(tmp) < 2 {
		return "", "", 0, errors.New("url is invalid")
	}
	pass := tmp[0]
	host := tmp[1]

	return host, pass, port, nil
}
