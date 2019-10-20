package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	conf "github.com/hiromaily/go-gin-wrapper/pkg/configs"
	"github.com/hiromaily/go-gin-wrapper/pkg/libs/fcgi"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/libs/ginsession"
	"github.com/hiromaily/go-gin-wrapper/pkg/routes"
	"github.com/hiromaily/golibs/auth/jwt"
	enc "github.com/hiromaily/golibs/cipher/encryption"
	mongo "github.com/hiromaily/golibs/db/mongodb"
	"github.com/hiromaily/golibs/db/mysql"
	fl "github.com/hiromaily/golibs/files"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/signal"
)

var (
	tomlPath = flag.String("f", "", "Toml file path")
	portNum  = flag.Int("P", 0, "Port of server")
)

func init() {
	//command-line
	flag.Parse()

	//cipher
	enc.NewCryptWithEnv()
}

func setupMain() {
	//conf
	initConf()

	//log
	lg.InitializeLog(lg.LogStatus(conf.GetConf().Server.Log.Level), lg.TimeShortFile,
		"[GOWEB]", conf.GetConf().Server.Log.Path, "hiromaily")

	//lg.Debugf("conf %#v\n", conf.GetConfInstance())
	lg.Debugf("[Environment] : %s\n", conf.GetConf().Environment)

	//auth settings
	initAuth()

	// debug mode
	if conf.GetConf().Environment == "local" {
		//signal
		go signal.StartSignal()
	}
	if conf.GetConf().Environment == "production" {
		//For release
		gin.SetMode(gin.ReleaseMode)
	}

	//Database settings
	initDatabase(0)
}

func initConf() {
	//config
	conf.New(*tomlPath, true)
}

func initAuth() {
	auth := conf.GetConf().API.JWT
	if auth.Mode == jwt.HMAC && auth.Secret != "" {
		jwt.InitSecretKey(auth.Secret)
	} else if auth.Mode == jwt.RSA && auth.PrivateKey != "" && auth.PublicKey != "" {
		err := jwt.InitKeys(auth.PrivateKey, auth.PublicKey)
		if err != nil {
			lg.Error(err)
			panic(err)
		}
	} else {
		jwt.InitEncrypted(jwt.HMAC)
		//lg.Debug("JWT Auth is not available because of toml settings.")
	}
}

// initialize Database
func initDatabase(testFlg uint8) {
	//if os.Getenv("HEROKU_FLG") == "1" {
	if conf.GetConf().Environment == "heroku" {
		//Heroku
		lg.Debug("HEROKU mode")

		//database
		host, dbname, user, pass, err := hrk.GetMySQLInfo("")
		//lg.Debugf("[HOST]%s  [Database]%s", host, dbname)
		//lg.Debugf("[User]%s  [Pass]%s", user, pass)

		if err != nil {
			lg.Error(err)
			panic(err)
		} else {
			dbInfo := conf.GetConf().MySQL
			dbInfo.Host = host
			dbInfo.DbName = dbname
			dbInfo.User = user
			dbInfo.Pass = pass
			dbInfo.Port = 3306
		}
	}

	//database
	if testFlg == 0 {
		dbInfo := conf.GetConf().MySQL
		mysql.New(dbInfo.Host, dbInfo.DbName, dbInfo.User, dbInfo.Pass, dbInfo.Port)
	} else {
		//For test
		dbInfo := conf.GetConf().MySQL.Test
		mysql.New(dbInfo.Host, dbInfo.DbName, dbInfo.User, dbInfo.Pass, dbInfo.Port)
	}
	mysql.GetDB().SetMaxIdleConns(50)

	//MongoDB
	initMongo()
}

func initMongo() {
	c := conf.GetConf().Mongo

	if conf.GetConf().Environment == "heroku" {
		host, dbname, user, pass, port, err := hrk.GetMongoInfo("")
		if err == nil {
			c.Host = host
			c.DbName = dbname
			c.User = user
			c.Pass = pass
			c.Port = uint16(port)
		}
	}

	mongo.New(c.Host, c.DbName, c.User, c.Pass, c.Port)
	if c.DbName != "" {
		//GetMongo().GetDB("hiromaily")
		mongo.GetMongo().GetDB(c.DbName)
	}
}

// initialize session
func initSession(r *gin.Engine) {
	red := conf.GetConf().Redis
	//if os.Getenv("HEROKU_FLG") == "1" {
	if conf.GetConf().Environment == "heroku" {
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
	//TODO:skip static files like (jpg, gif, png, js, css, woff)

	r.Use(gin.Logger())

	//r.Use(gin.Recovery())         //After GlobalRecover()
	r.Use(routes.GlobalRecover()) //It's called faster than [gin.Recovery()]

	//session
	initSession(r)

	//TODO:set ip to toml or redis server
	//check ip address to refuse specific IP Address
	//when using load balancer or reverse proxy, set specific IP
	r.Use(routes.RejectSpecificIP())

	//meta data for each rogic
	r.Use(routes.SetMetaData())

	//auto session(expire) update
	r.Use(routes.UpdateUserSession())
}

func getPort() (port int) {
	//For Heroku
	//if os.Getenv("PORT") != "" {
	if *portNum != 0 {
		//port = u.Atoi(os.Getenv("PORT"))
		port = *portNum
		conf.GetConf().Server.Port = port
	} else {
		port = conf.GetConf().Server.Port
	}
	lg.Debugf("port:%d", port)

	return
}

func loadTemplates(r *gin.Engine) {
	//http://stackoverflow.com/questions/25745701/parseglob-what-is-the-pattern-to-parse-all-templates-recursively-within-a-direc

	//r.LoadHTMLGlob("templates/*")
	//r.LoadHTMLGlob("templates/**/*")

	//It's impossible to call more than one. it was overwritten by final call.
	//r.LoadHTMLGlob(path + "templates/pages/**/*")
	//r.LoadHTMLGlob(path + "templates/components/*")

	//rootPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-gin-wrapper"
	rootPath := conf.GetConf().Server.Docs.Path

	ext := []string{"tmpl"}

	files1 := fl.GetFileList(rootPath+"/web/templates/pages", ext)
	files2 := fl.GetFileList(rootPath+"/web/templates/components", ext)
	files3 := fl.GetFileList(rootPath+"/web/templates/inner_js", ext)

	joined1 := append(files1, files2...)
	files := append(joined1, files3...)

	//tmpls := template.Must(template.ParseFiles(files...))
	//tmpls := template.Must(template.ParseFiles(files...)).Funcs(getTempFunc())
	tmpls := template.Must(template.New("").Funcs(getTempFunc()).ParseFiles(files...))
	//
	r.SetHTMLTemplate(tmpls)
}

// template FuncMap
func getTempFunc() template.FuncMap {
	//type FuncMap map[string]interface{}

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"strAry": func(ary []string, i int) string {
			return ary[i]
		},
		"dateFmt": func(t time.Time) string {
			//fmt := "August 17, 2016 9:51 pm"
			//fmt := "2006-01-02 15:04:05"
			//fmt := "Mon Jan _2 15:04:05 2006"
			fmt := "Mon Jan _2 15:04:05"
			return t.Format(fmt)
		},
	}
	return funcMap
}

func loadStaticFiles(r *gin.Engine) {
	//rootPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-gin-wrapper"
	rootPath := conf.GetConf().Server.Docs.Path

	//r.Static("/static", "/var/www")
	r.Static("/statics", rootPath+"/statics")
	r.Static("/assets", rootPath+"/statics/assets")
	r.Static("/favicon.ico", rootPath+"/statics/favicon.ico")
	r.Static("/swagger", rootPath+"/swagger/swagger-ui")

	// /when location of html as layer level is not top, be careful.
	//r.Static("/admin/assets", "statics/assets")
}

func run(r *gin.Engine) {
	port := getPort()
	if conf.GetConf().Proxy.Mode == 2 {
		//Proxy(Nginx) settings
		color.Red("[WARNING] running on fcgi mode.")
		lg.Info("running on fcgi mode.")
		fcgi.Run(r, fmt.Sprintf(":%d", port))
	} else {
		r.Run(fmt.Sprintf(":%d", port))
		//change to endless for Zero downtime restarts
		//endless.ListenAndServe(fmt.Sprintf(":%d", port), r)
	}
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
	loadTemplates(r)

	// Static
	loadStaticFiles(r)

	// Set router
	routes.SetURLOnHTTP(r)

	// Set Profiling
	if conf.GetConf().Develop.ProfileEnable {
		ginpprof.Wrapper(r)
	}

	// When Testing
	if testFlg == 1 {
		return r
	}

	// Run
	run(r)

	return r
}

// For TLS (work in progress)
func setHTTPSServer() {
	//gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	//Global middleware
	setMiddleWare(r)

	//templates
	r.LoadHTMLGlob("web/templates/**/*")

	//static
	//router.Static("/static", "/var/www")
	r.Static("/statics", "web/statics")
	r.Static("/assets", "web/statics/assets")

	//set router
	routes.SetURLOnHTTPS(r)

	// [HTTPS] TSL
	//r.RunTLS(addr string, certFile string, keyFile string)
}

// Creates a gin router with default middleware:
// logger and recovery (crash-free) middleware
func main() {
	setupMain()

	//HTTP
	setHTTPServer(0, "")

	//HTTPS
	//setHTTPSServer()
}
