package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ctls "github.com/hiromaily/go-gin-wrapper/pkg/server/controllers"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/cors"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/middlewares"
	//"github.com/hiromaily/go-gin-wrapper/controllers/chat"
)

// SetURLOnHTTP is for HTTP
func (s *Server) SetURLOnHTTP(r *gin.Engine) {
	// m := melody.New()

	ctl := ctls.NewController(
		s.dbModeler,
		// s.kvsStorager.CreateDBModel(),
		s.mongoModeler,
		s.conf.API.Header,
		s.conf.Auth,
		s.conf.API.CORS,
	)

	/******************************************************************************/
	/******** Return HTML *********************************************************/
	/******************************************************************************/
	//TODO:When http request method is POST, check referer in advance automatically.
	//-----------------------
	// Base(Top Level)
	//-----------------------

	r.GET("/", ctl.BaseIndexAction)
	// Redirect
	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	r.HEAD("/", func(c *gin.Context) {}) // For helth check

	// Login
	r.GET("/login", ctl.BaseLoginGetAction)
	r.POST("/login", middlewares.CheckHTTPRefererAndCSRF(s.conf.Server), ctl.BaseLoginPostAction)

	// Loout
	r.PUT("/logout", ctl.BaseLogoutPutAction)   // For Ajax
	r.POST("/logout", ctl.BaseLogoutPostAction) // HTML

	//-----------------------
	// News
	//-----------------------
	newsG := r.Group("/news")
	{
		// Top
		newsG.GET("/", ctl.NewsIndexAction)
	}
	// r.GET("/news/", news.NewsGetAction)

	//-----------------------
	// API List
	//-----------------------
	apiListG := r.Group("/apilist")
	{
		// Top
		apiListG.GET("/", ctl.APIListIndexAction)
		apiListG.GET("/index2", ctl.APIListIndex2Action)
	}

	//-----------------------
	// Chat (Work in progress)
	//-----------------------
	//chatG := r.Group("/chat")
	//{
	//	//Top
	//	chatG.GET("/", chat.IndexAction)
	//}
	//r.GET("/ws", func(c *gin.Context) {
	//	m.HandleRequest(c.Writer, c.Request)
	//})
	//m.HandleMessage(func(s *melody.Session, msg []byte) {
	//	m.Broadcast(msg)
	//})

	//-----------------------
	// OAuth2 Callback
	//-----------------------
	oauth2G := r.Group("/oauth2")
	{
		//--Google--
		//Sign in
		oauth2G.GET("/google/signin", ctl.OAuth2SignInGoogleAction)
		// Callback
		oauth2G.GET("/google/callback", ctl.OAuth2CallbackGoogleAction)
		//--Facebook--
		//Sign in
		oauth2G.GET("/facebook/signin", ctl.OAuth2SignInFacebookAction)
		// Callback
		oauth2G.GET("/facebook/callback", ctl.OAuth2CallbackFacebookAction)
	}

	//-----------------------
	// Account(MyPage)
	//-----------------------
	//After login
	accountsG := r.Group("/accounts")
	{
		// Top
		accountsG.GET("/", ctl.AccountIndexAction)
	}
	// r.GET("/accounts/", accounts.AccountsGetAction)

	//-----------------------
	// Admin [BasicAuth()]
	//-----------------------
	ba := s.conf.Server.BasicAuth
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		ba.User: ba.Pass,
	}))

	authorized.GET("/", ctl.AdminIndexAction)
	authorized.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	//-----------------------
	/* Error HTML */
	//-----------------------
	r.NoRoute(ctl.Error404Action)
	r.NoMethod(ctl.Error405Action)

	/******************************************************************************/
	/************ REST API (For Ajax) *********************************************/
	/******************************************************************************/
	//-----------------------
	// JWT
	//-----------------------
	jwt := r.Group("/api/jwt", middlewares.CheckHTTPHeader(s.conf.API))
	{
		jwt.POST("", ctl.APIJWTIndexPostAction) // jwt end point
	}

	//-----------------------
	// User
	//-----------------------
	//TODO: which is better to use CheckHttpHeader() for ajax request to check http header.
	// if it's used as middle ware like as below.
	//  r.Use(routes.CheckHttpHeader())
	//  it let us faster to develop instead of a bit less performance.
	handlers := []gin.HandlerFunc{middlewares.CheckHTTPHeader(s.conf.API)}
	// JWT
	if s.conf.API.JWT.Mode != 0 {
		handlers = append(handlers, middlewares.CheckJWT())
	}
	// CORS
	if s.conf.API.CORS.Enabled {
		handlers = append(handlers, middlewares.CheckCORS(s.conf.API.CORS))
	}

	// users := r.Group("/api/users", CheckHttpHeader(), CheckJWT())
	users := r.Group("/api/users", handlers...)
	{
		users.GET("", ctl.APIUserListGetAction)          // Get user list
		users.POST("", ctl.APIUserInsertPostAction)      // Register for new user
		users.GET("/id/:id", ctl.APIUserGetAction)       // Get specific user
		users.PUT("/id/:id", ctl.APIUserPutAction)       // Update specific user
		users.DELETE("/id/:id", ctl.APIUserDeleteAction) // Delete specific user
		// for unnecessary parameter, use *XXXX. e.g. /user/:name/*action

		// panic: path segment 'ids' conflicts with existing wildcard ':id' in path '/api/users/ids'
		users.GET("/ids", ctl.APIUserIDsGetAction) // Get user list

		// Accept CORS
		users.OPTIONS("", cors.SetHeader(s.conf.API.CORS))
	}

	// TODO:When user can use only method of GET and POST, X-HTTP-Method-Override header may be helpful.
	// Or use parameter `_method`
}
