package main

import (
	"flag"
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	"github.com/hiromaily/go-gin-wrapper/libs/fcgi"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	"github.com/hiromaily/go-gin-wrapper/routes"
	"github.com/hiromaily/golibs/db/mysql"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/signal"
	u "github.com/hiromaily/golibs/utils"
	"os"
)

var (
	tomlPath = flag.String("f", "", "Toml file path")
)

func init() {
	//command-line
	flag.Parse()

	//conf
	initConf()

	//log
	lg.InitializeLog(conf.GetConfInstance().Server.Log.Level, lg.LOG_OFF_COUNT, 0,
		"[GOWEB]", conf.GetConfInstance().Server.Log.Path)

	//lg.Debugf("conf %#v\n", conf.GetConfInstance())
	lg.Debugf("[Environment] : %s\n", conf.GetConfInstance().Environment)

	//Database settings
	//if os.Getenv("HEROKU_FLG") == "1" {
	if conf.GetConfInstance().Environment == "heroku" {
		//Heroku
		lg.Debug("HEROKU mode")

		//database
		host, dbname, user, pass, err := hrk.GetMySQLInfo("")
		lg.Debugf("[HOST]%s  [Database]%s", host, dbname)
		lg.Debugf("[User]%s  [Pass]%s", user, pass)

		if err != nil {
			panic(err)
		} else {
			mysql.New(host, dbname, user, pass, 3306)
			return
		}
	} else {
		//For Localhost, Docker, Production

		//database
		dbInfo := conf.GetConfInstance().MySQL
		mysql.New(dbInfo.Host, dbInfo.DbName, dbInfo.User, dbInfo.Pass, dbInfo.Port)
	}

	// debug mode
	if conf.GetConfInstance().Environment == "local" {
		//signal
		go signal.StartSignal()
	} else if conf.GetConfInstance().Environment == "production" {
		//For release
		gin.SetMode(gin.ReleaseMode)
	}
}

func initConf() {
	//config
	if *tomlPath != "" {
		conf.SetTomlPath(*tomlPath)
	} else {
		//default on localhost
		tomlPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-gin-wrapper/configs/settings.toml"
		conf.SetTomlPath(tomlPath)
	}
}

// initialize session
func initSession(r *gin.Engine) {
	red := conf.GetConfInstance().Redis
	//if os.Getenv("HEROKU_FLG") == "1" {
	if conf.GetConfInstance().Environment == "heroku" {
		host, pass, port, err := hrk.GetRedisInfo("")
		if err == nil && host != "" && port != 0 {
			red.Session = true
			red.Host = host
			red.Port = uint16(port)
			red.Pass = pass
		}
	}

	if red.Session && red.Host != "" && red.Port != 0 {
		lg.Debug("redis session start")
		sess.SetSession(r, fmt.Sprintf("%s:%d", red.Host, red.Port), red.Pass)
	} else {
		sess.SetSession(r, "", "")
	}
}

//Global middleware
func setMiddleWare(r *gin.Engine) {

	r.Use(gin.Logger())

	//r.Use(gin.Recovery())         //After GlobalRecover()
	r.Use(routes.GlobalRecover()) //It's called faster than [gin.Recovery()]

	//session
	initSession(r)

	//TODO:set ip to toml or redis server
	//check ip address to refuse specific IP Address
	//when using load balancer or reverse proxy, set specific IP
	r.Use(routes.RejectSpecificIp())

	//auto session(expire) update
	r.Use(routes.UpdateUserSession())

	//meta data for each rogic
	r.Use(routes.SetMetaData())

}

func getPort() (port int) {
	if os.Getenv("PORT") == "" {
		port = conf.GetConfInstance().Server.Port
	} else {
		port = u.Atoi(os.Getenv("PORT"))
		//conf.GetConfInstance().Server.Port = port
	}
	lg.Debugf("port:%d", port)

	return
}

func setHTTPServer(testFlg uint8, path string) *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	//r := gin.Default()
	/*
		func Default() *Engine {
			engine := New()
			engine.Use(Logger(), Recovery())
			return engine
		}
	*/

	// Global middleware
	setMiddleWare(r)

	// Templates
	//router.LoadHTMLGlob("templates/*")
	//r.LoadHTMLGlob("templates/**/*")
	r.LoadHTMLGlob(path + "templates/**/*")

	// Static
	//router.Static("/static", "/var/www")
	r.Static("/statics", "statics")
	r.Static("/assets", "statics/assets")
	//when location of html as layer level is not top, be careful.
	//r.Static("/admin/assets", "statics/assets")

	// Set router
	routes.SetHTTPUrls(r)

	// Set Profiling
	if conf.GetConfInstance().Profile.Enable {
		ginpprof.Wrapper(r)
	}

	// Test
	if testFlg == 1 {
		return r
	}

	// Run
	port := getPort()
	if conf.GetConfInstance().Proxy.Enable {
		//Proxy(Nginx) settings
		fcgi.Run(r, fmt.Sprintf(":%d", port))
	} else {
		r.Run(fmt.Sprintf(":%d", port))
	}
	return r
}

// For TLS (work in progress)
func setHTTPSServer() {
	//gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	//Global middleware
	setMiddleWare(r)

	//templates
	r.LoadHTMLGlob("templates/**/*")

	//static
	//router.Static("/static", "/var/www")
	r.Static("/statics", "statics")
	r.Static("/assets", "statics/assets")

	//set router
	routes.SetHTTPSUrls(r)

	// [HTTPS] TSL
	//r.RunTLS(addr string, certFile string, keyFile string)
}

// Creates a gin router with default middleware:
// logger and recovery (crash-free) middleware
func main() {
	//HTTP
	setHTTPServer(0, "")

	//HTTPS
	//setHTTPSServer()
}
