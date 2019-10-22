package server

import (
	"fmt"
	"html/template"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/configs"
	"github.com/hiromaily/go-gin-wrapper/pkg/libs/fcgi"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/libs/ginsession"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mongo"
	dbmodel "github.com/hiromaily/go-gin-wrapper/pkg/models/mysql"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/middlewares"
	"github.com/hiromaily/go-gin-wrapper/pkg/storager/session"
	fl "github.com/hiromaily/golibs/files"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
)

// ----------------------------------------------------------------------------
// Serverer interface
// ----------------------------------------------------------------------------

// Serverer is Serverer interface
type Serverer interface {
	Start() error
	Close()
}

// NewServerer is to return Serverer interface
func NewServerer(
	conf *configs.Config,
	port int,
	dbModeler dbmodel.DBModeler,
	mongoModeler mongomodel.MongoModeler,
	sessionStorage session.SessionStorager) Serverer {

	return NewServer(conf, port, dbModeler, mongoModeler, sessionStorage)
}

// ----------------------------------------------------------------------------
// Server
// ----------------------------------------------------------------------------

// Server is Server object
type Server struct {
	conf            *configs.Config
	port            int
	dbModeler       dbmodel.DBModeler
	mongoModeler    mongomodel.MongoModeler
	sessionStorager session.SessionStorager
	gin             *gin.Engine
}

// NewServer is to return server object
func NewServer(
	conf *configs.Config,
	port int,
	dbModeler dbmodel.DBModeler,
	mongoModeler mongomodel.MongoModeler,
	sessionStorager session.SessionStorager) *Server {

	if port == 0 {
		port = conf.Server.Port
	}

	srv := Server{
		conf:            conf,
		port:            port,
		dbModeler:       dbModeler,
		mongoModeler:    mongoModeler,
		sessionStorager: sessionStorager,
		gin:             gin.New(),
	}
	return &srv
}

// Start is to start server execution
func (s *Server) Start() error {
	if s.conf.Environment == "production" {
		//For release
		gin.SetMode(gin.ReleaseMode)
	}

	// Global middleware
	s.setMiddleWare()

	// Templates
	s.loadTemplates()

	// Static
	s.loadStaticFiles()

	// Set router (from urls.go)
	s.SetURLOnHTTP(s.gin)

	// Set Profiling
	if s.conf.Develop.ProfileEnable {
		ginpprof.Wrapper(s.gin)
	}

	// Run
	s.run()

	//return r
	return nil
}

// Close is to clean up middleware object
// TODO: not implemented yet
func (s *Server) Close() {
	//s.storager.Close()
}

//Global middleware
func (s *Server) setMiddleWare() {
	//TODO:skip static files like (jpg, gif, png, js, css, woff)

	s.gin.Use(gin.Logger())

	//r.Use(gin.Recovery())           //After GlobalRecover()
	s.gin.Use(middlewares.GlobalRecover()) //It's called faster than [gin.Recovery()]

	//session
	s.initSession()

	//TODO:set ip to toml or redis server
	//check ip address to refuse specific IP Address
	//when using load balancer or reverse proxy, set specific IP
	s.gin.Use(middlewares.RejectSpecificIP())

	//meta data for each rogic
	s.gin.Use(middlewares.SetMetaData())

	//auto session(expire) update
	s.gin.Use(middlewares.UpdateUserSession())
}

// TODO: it should be used as local object
func (s *Server) initSession() {
	red := s.conf.Redis
	//if os.Getenv("HEROKU_FLG") == "1" {
	if s.conf.Environment == "heroku" {
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
		sess.SetSession(s.gin, fmt.Sprintf("%s:%d", red.Host, red.Port), red.Pass)
	} else {
		sess.SetSession(s.gin, "", "")
	}
}
func (s *Server) loadTemplates() {
	//http://stackoverflow.com/questions/25745701/parseglob-what-is-the-pattern-to-parse-all-templates-recursively-within-a-direc

	//r.LoadHTMLGlob("templates/*")
	//r.LoadHTMLGlob("templates/**/*")

	//It's impossible to call more than one. it was overwritten by final call.
	//r.LoadHTMLGlob(path + "templates/pages/**/*")
	//r.LoadHTMLGlob(path + "templates/components/*")

	rootPath := s.conf.Server.Docs.Path

	ext := []string{"tmpl"}

	files1 := fl.GetFileList(rootPath+"/web/templates/pages", ext)
	files2 := fl.GetFileList(rootPath+"/web/templates/components", ext)
	files3 := fl.GetFileList(rootPath+"/web/templates/inner_js", ext)

	var files []string
	files = append(files, files1...)
	files = append(files, files2...)
	files = append(files, files3...)

	tmpls := template.Must(template.New("").Funcs(getTempFunc()).ParseFiles(files...))
	s.gin.SetHTMLTemplate(tmpls)
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

func (s *Server) loadStaticFiles() {
	//rootPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-gin-wrapper"
	rootPath := s.conf.Server.Docs.Path

	//r.Static("/static", "/var/www")
	s.gin.Static("/statics", rootPath+"/statics")
	s.gin.Static("/assets", rootPath+"/statics/assets")
	s.gin.Static("/favicon.ico", rootPath+"/statics/favicon.ico")
	s.gin.Static("/swagger", rootPath+"/swagger/swagger-ui")

	// /when location of html as layer level is not top, be careful.
	//r.Static("/admin/assets", "statics/assets")
}

func (s *Server) run() {
	if s.conf.Proxy.Mode == 2 {
		//Proxy(Nginx) settings
		color.Red("[WARNING] running on fcgi mode.")
		lg.Info("running on fcgi mode.")
		fcgi.Run(s.gin, fmt.Sprintf(":%d", s.port))
	} else {
		s.gin.Run(fmt.Sprintf(":%d", s.port))
		//change to endless for Zero downtime restarts
		//endless.ListenAndServe(fmt.Sprintf(":%d", port), r)
	}
}
