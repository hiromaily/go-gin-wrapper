package routes

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	"github.com/hiromaily/go-gin-wrapper/controllers/accounts"
	"github.com/hiromaily/go-gin-wrapper/controllers/admins"
	"github.com/hiromaily/go-gin-wrapper/controllers/api/jwt"
	us "github.com/hiromaily/go-gin-wrapper/controllers/api/users"
	"github.com/hiromaily/go-gin-wrapper/controllers/bases"
	"github.com/hiromaily/go-gin-wrapper/controllers/errors"
	"github.com/hiromaily/go-gin-wrapper/controllers/news"
	oauth "github.com/hiromaily/go-gin-wrapper/controllers/oauth2"
	//ba "github.com/hiromaily/go-gin-wrapper/libs/basicauth"
	"github.com/hiromaily/go-gin-wrapper/controllers/apilist"
	//"github.com/hiromaily/go-gin-wrapper/controllers/chat"
	//"github.com/olahol/melody"
	"net/http"
)

// RefererURLs key->request url, value->refer url
var RefererURLs = map[string]string{
	"/login": "login",
	"/user":  "user",
}

// SetURLOnHTTP is for HTTP
func SetURLOnHTTP(r *gin.Engine) {
	//m := melody.New()

	/******************************************************************************/
	/******** Return HTML *********************************************************/
	/******************************************************************************/
	//TODO:When http request method is POST, check referer in advance automatically.
	//-----------------------
	// Base(Top Level)
	//-----------------------
	r.GET("/", bases.IndexAction)
	//Redirect
	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	r.HEAD("/", func(c *gin.Context) {}) //For helth check

	//Login
	r.GET("/login", bases.LoginGetAction)
	r.POST("/login", CheckHTTPRefererAndCSRF(), bases.LoginPostAction)

	//Loout
	r.PUT("/logout", bases.LogoutPutAction)   //For Ajax
	r.POST("/logout", bases.LogoutPostAction) //HTML

	//-----------------------
	// News
	//-----------------------
	newsG := r.Group("/news")
	{
		//Top
		newsG.GET("/", news.IndexAction)
	}
	//r.GET("/news/", news.NewsGetAction)

	//-----------------------
	// API List
	//-----------------------
	apiListG := r.Group("/apilist")
	{
		//Top
		apiListG.GET("/", apilist.IndexAction)
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
		oauth2G.GET("/google/signin", oauth.SignInGoogleAction)
		//Callback
		oauth2G.GET("/google/callback", oauth.CallbackGoogleAction)
		//--Facebook--
		//Sign in
		oauth2G.GET("/facebook/signin", oauth.SignInFacebookAction)
		//Callback
		oauth2G.GET("/facebook/callback", oauth.CallbackFacebookAction)
	}

	//-----------------------
	// Account(MyPage)
	//-----------------------
	//After login
	accountsG := r.Group("/accounts")
	{
		//Top
		accountsG.GET("/", accounts.IndexAction)
	}
	//r.GET("/accounts/", accounts.AccountsGetAction)

	//-----------------------
	// Admin [BasicAuth()]
	//-----------------------
	ba := conf.GetConf().Server.BasicAuth
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		ba.User: ba.Pass,
	}))

	authorized.GET("/", admins.IndexAction)
	authorized.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	//-----------------------
	/* Error HTML */
	//-----------------------
	r.NoRoute(errors.Error404Action)
	r.NoMethod(errors.Error405Action)

	/******************************************************************************/
	/************ REST API (For Ajax) *********************************************/
	/******************************************************************************/
	//-----------------------
	// JWT
	//-----------------------
	jw := r.Group("/api/jwt", CheckHTTPHeader())
	{
		jw.POST("", jwt.IndexPostAction) //jwt end point
	}

	//-----------------------
	// User
	//-----------------------
	//TODO: which is better to use CheckHttpHeader() for ajax request to check http header.
	// if it's used as middle ware like as below.
	//  r.Use(routes.CheckHttpHeader())
	//  it let us faster to develop instead of a bit less performance.
	var handlers = []gin.HandlerFunc{CheckHTTPHeader()}
	if conf.GetConf().Auth.JWT.Mode != 0 {
		handlers = append(handlers, CheckJWT())
	}

	//users := r.Group("/api/users", CheckHttpHeader(), CheckJWT())
	users := r.Group("/api/users", handlers...)
	{
		users.GET("", us.ListGetAction)       //Get user list
		users.POST("", us.InsertPostAction)   //Register for new user
		users.GET("/:id", us.GetAction)       //Get specific user
		users.PUT("/:id", us.PutAction)       //Update specific user
		users.DELETE("/:id", us.DeleteAction) //Delete specific user
		//for unnecessary parameter, use *XXXX. e.g. /user/:name/*action
	}
	//TODO:When user can use only method of GET and POST, X-HTTP-Method-Override header may be helpful.
	//Or use parameter `_method`
}

// SetURLOnHTTPS is for HTTPS
// TODO: it may not be necessary
func SetURLOnHTTPS(r *gin.Engine) {

}
