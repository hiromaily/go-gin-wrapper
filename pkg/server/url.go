package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/cors"
)

// setRouter is for HTTP
func (s *server) setRouter(r *gin.Engine) {
	s.logger.Info("server setRouter()")
	/******************************************************************************/
	/******** Return HTML *********************************************************/
	/******************************************************************************/
	//TODO:When http request method is POST, check referer in advance automatically.
	//-----------------------
	// Base(Top Level)
	//-----------------------

	r.GET("/", s.controller.BaseIndexAction)
	// Redirect
	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	r.HEAD("/", func(c *gin.Context) {}) // For helth check

	// Login
	r.GET("/login", s.controller.BaseLoginGetAction)
	r.POST("/login", s.middleware.CheckHTTPRefererAndCSRF(), s.controller.BaseLoginPostAction)

	// Loout
	r.PUT("/logout", s.controller.BaseLogoutPutAction)   // For Ajax
	r.POST("/logout", s.controller.BaseLogoutPostAction) // HTML

	//-----------------------
	// News
	//-----------------------
	//newsG := r.Group("/news")
	//{
	//	// Top
	//	newsG.GET("/", s.controller.NewsIndexAction)
	//}

	//-----------------------
	// API List
	//-----------------------
	apiListG := r.Group("/apilist")
	{
		// Top
		apiListG.GET("/", s.controller.APIListIndexAction)
		apiListG.GET("/index2", s.controller.APIListIndex2Action)
	}

	//-----------------------
	// OAuth2 Callback
	//-----------------------
	oauth2G := r.Group("/oauth2")
	{
		//--Google--
		//Sign in
		oauth2G.GET("/google/signin", s.controller.OAuth2SignInGoogleAction)
		// Callback
		oauth2G.GET("/google/callback", s.controller.OAuth2CallbackGoogleAction)
		//--Facebook--
		//Sign in
		oauth2G.GET("/facebook/signin", s.controller.OAuth2SignInFacebookAction)
		// Callback
		oauth2G.GET("/facebook/callback", s.controller.OAuth2CallbackFacebookAction)
	}

	//-----------------------
	// Account(MyPage)
	//-----------------------
	//After login
	accountsG := r.Group("/accounts")
	{
		// Top
		accountsG.GET("/", s.controller.AccountIndexAction)
	}
	// r.GET("/accounts/", accounts.AccountsGetAction)

	//-----------------------
	// Admin [BasicAuth()]
	//-----------------------
	basicAuth := s.serverConf.BasicAuth
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		basicAuth.User: basicAuth.Pass,
	}))

	authorized.GET("/", s.controller.AdminIndexAction)
	authorized.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	//-----------------------
	/* Error HTML */
	//-----------------------
	r.NoRoute(s.controller.Error404Action)
	r.NoMethod(s.controller.Error405Action)

	/******************************************************************************/
	/************ REST API (For Ajax) *********************************************/
	/******************************************************************************/
	//-----------------------
	// JWT
	//-----------------------
	jwt := r.Group("/api/jwt", s.middleware.CheckHTTPHeader())
	{
		jwt.POST("", s.controller.APIJWTIndexPostAction) // jwt end point
	}

	//-----------------------
	// User
	//-----------------------
	//TODO: which is better to use CheckHttpHeader() for ajax request to check http header.
	// if it's used as middle ware like as below.
	//  r.Use(routes.CheckHttpHeader())
	//  it let us faster to develop instead of a bit less performance.
	handlers := []gin.HandlerFunc{s.middleware.CheckHTTPHeader()}
	// JWT
	if s.apiConf.JWT.Mode != 0 {
		handlers = append(handlers, s.middleware.CheckJWT())
	}
	// CORS
	if s.apiConf.CORS.Enabled {
		handlers = append(handlers, s.middleware.CheckCORS())
	}

	// users := r.Group("/api/users", CheckHttpHeader(), CheckJWT())
	users := r.Group("/api/users", handlers...)
	{
		users.GET("", s.controller.APIUserListGetAction)          // Get user list
		users.POST("", s.controller.APIUserInsertPostAction)      // Register for new user
		users.GET("/id/:id", s.controller.APIUserGetAction)       // Get specific user
		users.PUT("/id/:id", s.controller.APIUserPutAction)       // Update specific user
		users.DELETE("/id/:id", s.controller.APIUserDeleteAction) // Delete specific user
		// for unnecessary parameter, use *XXXX. e.g. /user/:name/*action

		// panic: path segment 'ids' conflicts with existing wildcard ':id' in path '/api/users/ids'
		users.GET("/ids", s.controller.APIUserIDsGetAction) // Get user list

		// Accept CORS
		users.OPTIONS("", cors.SetHeader(s.logger, s.apiConf.CORS))
	}

	// TODO:When user can use only method of GET and POST, X-HTTP-Method-Override header may be helpful.
	// Or use parameter `_method`
}
