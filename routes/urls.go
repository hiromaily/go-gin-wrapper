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

//RefererUrl key->request url, value->refer url
var RefererUrls = map[string]string{
	"/login": "login",
	"/user":  "user",
}

//For HTTP
func SetHTTPUrls(r *gin.Engine) {
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
	r.POST("/login", CheckHttpRefererAndCSRF(), bases.LoginPostAction)

	//Loout
	r.PUT("/logout", bases.LogoutPutAction)   //For Ajax
	r.POST("/logout", bases.LogoutPostAction) //HTML

	//-----------------------
	// News
	//-----------------------
	newsG := r.Group("/news")
	{
		//Top
		newsG.GET("/", news.NewsGetAction)
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
		accountsG.GET("/", accounts.AccountsGetAction)
	}
	//r.GET("/accounts/", accounts.AccountsGetAction)

	//-----------------------
	// Admin [BasicAuth()]
	//-----------------------
	ba := conf.GetConf().Server.BasicAuth
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		ba.User: ba.Pass,
	}))
	//authorized := r.Group("/admin", ba.BasicAuth(ba.Accounts{
	//	"web": "test",
	//}))
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
	jw := r.Group("/api/jwt", CheckHttpHeader())
	{
		jw.POST("", jwt.IndexAction) //jwt end point
	}

	//-----------------------
	// User
	//-----------------------
	//TODO: which is better to use CheckHttpHeader() for ajax request to check http header.
	// if it's used as middle ware like as below.
	//  r.Use(routes.CheckHttpHeader())
	//  it let us faster to develop instead of a bit less performance.
	var handlers []gin.HandlerFunc = []gin.HandlerFunc{CheckHttpHeader()}
	if conf.GetConf().Auth.Jwt.Mode != 0 {
		handlers = append(handlers, CheckJWT())
	}

	//users := r.Group("/api/users", CheckHttpHeader(), CheckJWT())
	users := r.Group("/api/users", handlers...)
	{
		users.GET("", us.UsersListGetAction)      //Get user list
		users.POST("", us.UserPostAction)         //Register for new user
		users.GET("/:id", us.UserGetAction)       //Get specific user
		users.PUT("/:id", us.UserPutAction)       //Update specific user
		users.DELETE("/:id", us.UserDeleteAction) //Delete specific user
		//for unnecessary parameter, use *XXXX. e.g. /user/:name/*action
	}
	//TODO:When user can use only method of GET and POST, X-HTTP-Method-Override header may be helpful.
	//Or use parameter `_method`
}

//For HTTPS
func SetHTTPSUrls(r *gin.Engine) {

}
