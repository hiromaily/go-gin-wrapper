package server

import (
	"fmt"
	"html/template"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/dir"
	"github.com/hiromaily/go-gin-wrapper/pkg/repository"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/controller"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/fcgi"
	sess "github.com/hiromaily/go-gin-wrapper/pkg/server/ginsession"
	hrk "github.com/hiromaily/golibs/heroku"
)

// Server interface
type Server interface {
	Start() (*gin.Engine, error)
	Close()
}

// server object
type server struct {
	gin        *gin.Engine
	port       int
	middleware Middlewarer
	controller *controller.Controller // TODO: interface
	logger     *zap.Logger
	userRepo   repository.UserRepositorier

	serverConf  *config.ServerConfig
	proxyConf   *config.ProxyConfig
	apiConf     *config.APIConfig
	redisConf   *config.RedisConfig
	developConf *config.DevelopConfig

	isTestMode bool
}

// NewServer returns Server interface
func NewServer(
	gin *gin.Engine,
	port int,
	middleware Middlewarer,
	controller *controller.Controller,
	logger *zap.Logger,
	userRepo repository.UserRepositorier,
	conf *config.Config,
	isTestMode bool,
) Server {
	if port == 0 {
		port = conf.Server.Port
	}

	return &server{
		gin:         gin,
		port:        port,
		middleware:  middleware,
		controller:  controller,
		logger:      logger,
		userRepo:    userRepo,
		serverConf:  conf.Server,
		proxyConf:   conf.Proxy,
		apiConf:     conf.API,
		redisConf:   conf.Redis,
		developConf: conf.Develop,
		isTestMode:  isTestMode,
	}
}

// Start is to start server execution
func (s *server) Start() (*gin.Engine, error) {
	if s.serverConf.IsRelease {
		// For release
		gin.SetMode(gin.ReleaseMode)
	}

	// Global middleware
	s.setMiddleWare()

	// Templates
	s.loadTemplates()

	// Static
	s.loadStaticFiles()

	// Set router (from url.go)
	s.SetURLOnHTTP(s.gin)

	// Set Profiling
	if s.developConf.ProfileEnable {
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
func (s *server) Close() {
	// s.storager.Close()
}

// Global middleware
func (s *server) setMiddleWare() {
	// TODO:skip static files like (jpg, gif, png, js, css, woff)

	s.gin.Use(gin.Logger())

	// r.Use(gin.Recovery())  //After GlobalRecover()
	s.gin.Use(s.middleware.GlobalRecover()) // It's called faster than [gin.Recovery()]

	// session
	s.initSession()

	// TODO:set ip to toml or redis server
	// check ip address to refuse specific IP Address
	// when using load balancer or reverse proxy, set specific IP
	s.gin.Use(s.middleware.RejectSpecificIP())

	// meta data for each rogic
	s.gin.Use(s.middleware.SetMetaData())

	// auto session(expire) update
	s.gin.Use(s.middleware.UpdateUserSession())
}

func (s *server) initSession() {
	s.herokuRedisSetting()
	if s.redisConf.IsSession && s.redisConf.Host != "" && s.redisConf.Port != 0 {
		s.logger.Debug("initSession: redis session start")
		sess.SetSession(s.gin, s.logger, fmt.Sprintf("%s:%d", s.redisConf.Host, s.redisConf.Port), s.redisConf.Pass, s.serverConf.Session)
	} else {
		sess.SetSession(s.gin, s.logger, "", "", s.serverConf.Session)
	}
}

func (s *server) herokuRedisSetting() {
	if s.redisConf.IsHeroku {
		host, pass, port, err := hrk.GetRedisInfo("")
		if err == nil && host != "" && port != 0 {
			s.redisConf.IsSession = true
			s.redisConf.Host = host
			s.redisConf.Port = uint16(port)
			s.redisConf.Pass = pass
		}
	}
}

func (s *server) loadTemplates() {
	// http://stackoverflow.com/questions/25745701/parseglob-what-is-the-pattern-to-parse-all-templates-recursively-within-a-direc

	// r.LoadHTMLGlob("templates/*")
	// r.LoadHTMLGlob("templates/**/*")

	// It's impossible to call more than one. it was overwritten by final call.
	// r.LoadHTMLGlob(path + "templates/pages/**/*")
	// r.LoadHTMLGlob(path + "templates/components/*")

	rootPath := s.serverConf.Docs.Path

	ext := []string{"tmpl"}

	files1 := dir.GetFileList(rootPath+"/web/templates/pages", ext)
	files2 := dir.GetFileList(rootPath+"/web/templates/components", ext)
	files3 := dir.GetFileList(rootPath+"/web/templates/inner_js", ext)

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

func (s *server) loadStaticFiles() {
	rootPath := s.serverConf.Docs.Path

	// r.Static("/static", "/var/www")
	s.gin.Static("/statics", rootPath+"/web/statics")
	s.gin.Static("/assets", rootPath+"/web/statics/assets")
	s.gin.Static("/favicon.ico", rootPath+"/web/statics/favicon.ico")
	s.gin.Static("/swagger", rootPath+"/web/swagger/swagger-ui")
}

func (s *server) run() error {
	if s.proxyConf.Mode == 2 {
		// Proxy(Nginx) settings
		color.Red("[WARNING] running on fcgi mode.")
		s.logger.Info("running on fcgi mode.")
		return fcgi.Run(s.gin, fmt.Sprintf(":%d", s.port))
	}
	return s.gin.Run(fmt.Sprintf(":%d", s.port))
}
