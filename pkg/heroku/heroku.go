package heroku

import (
	"fmt"
	"os"
	"strings"

	str "github.com/hiromaily/go-gin-wrapper/pkg/strings"
)

// GetRedisInfo is to get Redis information
func GetRedisInfo(url string) (host, pass string, port int, err error) {
	if url == "" {
		url = os.Getenv("REDIS_URL")
		if url == "" {
			err = fmt.Errorf("REDIS_URL was not found")
			return
		}
	}
	//redis://h:xxx@xxx:20819
	//<password>@<hostname>:<port>
	tmp := strings.Split(url, "//")
	tmp = strings.Split(tmp[1], ":")
	port = str.Atoi(tmp[2])

	tmp = strings.Split(tmp[1], "@")
	pass = tmp[0]

	tmp = strings.Split(tmp[1], ":")
	host = tmp[0]

	return
}
