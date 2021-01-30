package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// setRouter is mapper for request path and handler
func (s *server) setRouter(r *gin.Engine) {
	s.logger.Info("server setRouter()")

	// for html response
	s.setBaseRouter(r)
	s.setAPIListRouter(r)
	s.setOAuth2Router(r)
	s.setAccountRouter(r)
	s.setAdminRouter(r)
	s.setErrorRouter(r)

	// for API use
	s.setJWTRouter(r)
	s.setUserRouter(r)

	// TODO: X-HTTP-Method-Override header may be helpful when user can use only `GET`, `POST` method of GET and POST
	// Or use parameter `_method`
}

// base router (top level)
func (s *server) setBaseRouter(r *gin.Engine) {
	// top level
	r.GET("/", s.controller.BaseIndexAction)
	r.GET("/index", func(ctx *gin.Context) { // redirect
		ctx.Redirect(http.StatusMovedPermanently, "/")
	})
	// r.HEAD("/", func(ctx *gin.Context) {}) // health check
	// FIXME: helath check must have specific endpoint for that.

	// login
	r.GET("/login", s.controller.BaseLoginGetAction)
	r.POST("/login",
		s.middleware.CheckHTTPReferer(),
		s.middleware.CheckCSRF(), // FIXME: set proper timing
		s.controller.BaseLoginPostAction,
	)

	// logout
	r.PUT("/logout", s.controller.BaseLogoutPutAction)   // for ajax
	r.POST("/logout", s.controller.BaseLogoutPostAction) // html
}

// API list router
func (s *server) setAPIListRouter(r *gin.Engine) {
	apiListG := r.Group("/apilist")

	apiListG.GET("/", s.controller.APIListIndexAction)
	apiListG.GET("/index2", s.controller.APIListIndex2Action)
}

// OAuth2 router
func (s *server) setOAuth2Router(r *gin.Engine) {
	oauth2G := r.Group("/oauth2")

	// google sign in
	oauth2G.GET("/google/signin", s.controller.OAuth2SignInGoogleAction)
	// google callback
	oauth2G.GET("/google/callback", s.controller.OAuth2CallbackGoogleAction)
	// facebook sign in
	oauth2G.GET("/facebook/signin", s.controller.OAuth2SignInFacebookAction)
	// facebook callback
	oauth2G.GET("/facebook/callback", s.controller.OAuth2CallbackFacebookAction)
}

// account router (my page)
// - used after login
func (s *server) setAccountRouter(r *gin.Engine) {
	accountsG := r.Group("/accounts")

	accountsG.GET("/", s.controller.AccountIndexAction)
}

// admin router with basic auth
func (s *server) setAdminRouter(r *gin.Engine) {
	basicAuth := s.serverConf.BasicAuth
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		basicAuth.User: basicAuth.Pass,
	}))

	authorized.GET("/", s.controller.AdminIndexAction)
	authorized.GET("/index", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/")
	})
}

// error router
func (s *server) setErrorRouter(r *gin.Engine) {
	r.NoRoute(s.controller.Error404Action)
	r.NoMethod(s.controller.Error405Action)
}

// JWT API router
func (s *server) setJWTRouter(r *gin.Engine) {
	jwt := r.Group("/api/jwts",
		s.middleware.CheckHTTPHeader(),
		s.middleware.SetResponseHeader(),
	)

	jwt.POST("", s.controller.APIJWTIndexPostAction) // jwt end point
}

// user API router
func (s *server) setUserRouter(r *gin.Engine) {
	// additional middleware handler
	preHandlers := []gin.HandlerFunc{s.middleware.CheckHTTPHeader()}
	preHandlers = append(preHandlers, s.middleware.CheckJWT())
	if s.apiConf.CORS.Enabled {
		preHandlers = append(preHandlers, s.middleware.CheckCORS())
	}
	preHandlers = append(preHandlers, s.middleware.SetResponseHeader())
	preHandlers = append(preHandlers, s.middleware.SetCORSHeader())

	users := r.Group("/api/users", preHandlers...)

	users.GET("", s.controller.APIUserListGetAction)
	users.POST("", s.controller.APIUserInsertPostAction)
	users.GET("/id/:id", s.controller.APIUserGetAction)
	users.PUT("/id/:id", s.controller.APIUserPutAction)
	users.DELETE("/id/:id", s.controller.APIUserDeleteAction)

	users.GET("/ids", s.controller.APIUserIDsGetAction) // get user list

	// accept CORS
	// FIXME: how to handle that only header is enough to run
	// users.OPTIONS("", s.middleware.SetCORSHeader())
}
