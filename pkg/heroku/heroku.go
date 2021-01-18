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
	tmp = strings.Split(tmp[1], ":")
	port := str.Atoi(tmp[2])

	tmp = strings.Split(tmp[1], "@")
	pass := tmp[0]

	tmp = strings.Split(tmp[1], ":")
	host := tmp[0]

	return host, pass, port, nil
}
