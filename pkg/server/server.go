package server

import (
	"fmt"
	"html/template"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	mongomodel "github.com/hiromaily/go-gin-wrapper/pkg/model/mongo"
	dbmodel "github.com/hiromaily/go-gin-wrapper/pkg/model/mysql"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/fcgi"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/middlewares"
	fl "github.com/hiromaily/golibs/files"
	hrk "github.com/hiromaily/golibs/heroku"
	lg "github.com/hiromaily/golibs/log"
)

// ----------------------------------------------------------------------------
// Serverer interface
// ----------------------------------------------------------------------------

// Serverer is Serverer interface
type Serverer interface {
	Start() (*gin.Engine, error)
	Close()
}

// NewServerer is to return Serverer interface
func NewServerer(
	isTestMode bool,
	conf *config.Config,
	port int,
	dbModeler dbmodel.DBModeler,
	mongoModeler mongomodel.MongoModeler) Serverer {
	return NewServer(isTestMode, conf, port, dbModeler, mongoModeler)
}

// ----------------------------------------------------------------------------
// Server
// ----------------------------------------------------------------------------

// Server is Server object
type Server struct {
	isTestMode   bool
	conf         *config.Config
	port         int
	dbModeler    dbmodel.DBModeler
	mongoModeler mongomodel.MongoModeler
	gin          *gin.Engine
}

// NewServer is to return server object
func NewServer(
	isTestMode bool,
	conf *config.Config,
	port int,
	dbModeler dbmodel.DBModeler,
	mongoModeler mongomodel.MongoModeler) *Server {
	if port == 0 {
		port = conf.Server.Port
	}

	srv := Server{
		isTestMode:   isTestMode,
		conf:         conf,
		port:         port,
		dbModeler:    dbModeler,
		mongoModeler: mongoModeler,
		gin:          gin.New(),
	}
	return &srv
}

// Start is to start server execution
func (s *Server) Start() (*gin.Engine, error) {
	if s.conf.Environment == "production" {
		// For release
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

	if s.isTestMode {
		return s.gin, nil
	}

	// Run
	err := s.run()
	return nil, err
}

// Close is to clean up middleware object
// TODO: not implemented yet
func (s *Server) Close() {
	// s.storager.Close()
}

// Global middleware
func (s *Server) setMiddleWare() {
	// TODO:skip static files like (jpg, gif, png, js, css, woff)

	s.gin.Use(gin.Logger())

	// r.Use(gin.Recovery())           //After GlobalRecover()
	s.gin.Use(middlewares.GlobalRecover(s.conf.Develop)) // It's called faster than [gin.Recovery()]

	// session
	s.initSession()

	// TODO:set ip to toml or redis server
	// check ip address to refuse specific IP Address
	// when using load balancer or reverse proxy, set specific IP
	s.gin.Use(middlewares.RejectSpecificIP(s.conf.Proxy))

	// meta data for each rogic
	s.gin.Use(middlewares.SetMetaData())

	// auto session(expire) update
	s.gin.Use(middlewares.UpdateUserSession())
}

func (s *Server) initSession() {
	red := s.conf.Redis
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
		sess.SetSession(s.gin, fmt.Sprintf("%s:%d", red.Host, red.Port), red.Pass, s.conf.Server.Session)
	} else {
		sess.SetSession(s.gin, "", "", s.conf.Server.Session)
	}
}

func (s *Server) loadTemplates() {
	// http://stackoverflow.com/questions/25745701/parseglob-what-is-the-pattern-to-parse-all-templates-recursively-within-a-direc

	// r.LoadHTMLGlob("templates/*")
	// r.LoadHTMLGlob("templates/**/*")

	// It's impossible to call more than one. it was overwritten by final call.
	// r.LoadHTMLGlob(path + "templates/pages/**/*")
	// r.LoadHTMLGlob(path + "templates/components/*")

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
	// type FuncMap map[string]interface{}

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
			// fmt := "August 17, 2016 9:51 pm"
			// fmt := "2006-01-02 15:04:05"
			// fmt := "Mon Jan _2 15:04:05 2006"
			fmt := "Mon Jan _2 15:04:05"
			return t.Format(fmt)
		},
	}
	return funcMap
}

func (s *Server) loadStaticFiles() {
	rootPath := s.conf.Server.Docs.Path

	// r.Static("/static", "/var/www")
	s.gin.Static("/statics", rootPath+"/web/statics")
	s.gin.Static("/assets", rootPath+"/web/statics/assets")
	s.gin.Static("/favicon.ico", rootPath+"/web/statics/favicon.ico")
	s.gin.Static("/swagger", rootPath+"/web/swagger/swagger-ui")
}

func (s *Server) run() error {
	if s.conf.Proxy.Mode == 2 {
		// Proxy(Nginx) settings
		color.Red("[WARNING] running on fcgi mode.")
		lg.Info("running on fcgi mode.")
		return fcgi.Run(s.gin, fmt.Sprintf(":%d", s.port))
	}
	return s.gin.Run(fmt.Sprintf(":%d", s.port))
}
