package heroku

import (
	"fmt"
	u "github.com/hiromaily/golibs/utils"
	"os"
	"strings"
)

func GetMySQLInfo(url string) (host, dbname, user, pass string, err error) {
	//CLEARDB_DATABASE_URL: mysql://be2ebea7cda583:49eef93c@us-cdbr-iron-east-04.cleardb.net/heroku_aa95a7f43af0311?reconnect=true
	if url == "" {
		url = os.Getenv("CLEARDB_DATABASE_URL")
		if url == "" {
			err = fmt.Errorf("CLEARDB_DATABASE_URL was not found.")
			return
		}
	}
	//_, err = fmt.Sscanf(url, "mysql://%s:%s@%s/%s?reconnect=true", &user, &pass, &host, &dbname)
	tmp := strings.Split(url, "//")
	tmp = strings.Split(tmp[1], ":")
	user = tmp[0]
	tmp = strings.Split(tmp[1], "@")
	pass = tmp[0]
	tmp = strings.Split(tmp[1], "/")
	host = tmp[0]
	dbname = strings.Split(tmp[1], "?")[0]

	return
}

func GetRedisInfo(url string) (host, pass string, port int, err error) {
	if url == "" {
		url = os.Getenv("REDIS_URL")
		if url == "" {
			err = fmt.Errorf("REDIS_URL was not found.")
			return
		}
	}
	//redis://h:pf217irr4gts39d29o0012bghsi@ec2-50-19-83-130.compute-1.amazonaws.com:20819
	//<password>@<hostname>:<port>
	tmp := strings.Split(url, "//")
	tmp = strings.Split(tmp[1], ":")
	port = u.Atoi(tmp[2])

	tmp = strings.Split(tmp[1], "@")
	pass = tmp[0]

	tmp = strings.Split(tmp[1], ":")
	host = tmp[0]

	return
}

func GetMongoInfo(url string) (host, dbname, user, pass string, port int, err error) {
	//MONGODB_URI: mongodb://heroku_7lbnd77m:7r8f631nv2idt0fhj9ok9714j9@ds161495.mlab.com:61495/heroku_7lbnd77m
	if url == "" {
		url = os.Getenv("MONGODB_URI")
		if url == "" {
			err = fmt.Errorf("MONGODB_URI was not found.")
			return
		}
	}
	//_, err = fmt.Sscanf(url, "mobngod://%s:%s@%s/%s?reconnect=true", &user, &pass, &host, &dbname)
	tmp := strings.Split(url, "//")
	tmp = strings.Split(tmp[1], ":")
	user = tmp[0]

	tmp2 := strings.Split(tmp[2], "/")
	port = u.Atoi(tmp2[0])
	dbname = tmp2[1]

	tmp = strings.Split(tmp[1], "@")
	pass = tmp[0]
	host = tmp[1]

	return
}
