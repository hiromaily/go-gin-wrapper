package main

import (
	"github.com/gin-gonic/gin"
	//"reflect"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/mysql"
	"github.com/hiromaily/golibs/signal"
	conf "github.com/hiromaily/web/configs"
	"github.com/hiromaily/web/routes"
)

func init() {
	//log
	//lg.InitLog(1, 0, "./logs/")
	lg.InitLog(1, 0, "")

	//config
	lg.Debugf("conf %#v\n", conf.GetConfInstance())
	lg.Debugf("conf environment : %s\n", conf.GetConfInstance().Environment)
	lg.Debugf("conf server host :%s\n", conf.GetConfInstance().Server.Host)
	lg.Debugf("conf mysql host :%s\n", conf.GetConfInstance().MySQL.Host)

	//database
	dbInfo := conf.GetConfInstance().MySQL
	mysql.New(dbInfo.Host, dbInfo.DbName, dbInfo.User, dbInfo.Pass, dbInfo.Port)

	if conf.GetConfInstance().Environment == "local" {
		//signal
		go signal.StartSignal()
	}
}

//Global middleware
func setMiddleWare(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//session
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("ginsession", store))
}

func setWebServer() {
	//gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	//fmt.Println(reflect.TypeOf(router))

	//Global middleware
	setMiddleWare(r)

	//templates
	//router.LoadHTMLGlob("templates/*")
	r.LoadHTMLGlob("templates/**/*")

	//static
	r.Static("/statics/", "statics")

	//set router
	routes.SetUrls(r)

	// Listen and server on 0.0.0.0:8080
	//r.Run(":9999")
	r.Run(fmt.Sprintf(":%d", conf.GetConfInstance().Server.Port))
}

// Creates a gin router with default middleware:
// logger and recovery (crash-free) middleware
func main() {
	//
	setWebServer()
}
