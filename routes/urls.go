package routes

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hiromaily/go-gin-wrapper/controllers/accounts"
	"github.com/hiromaily/go-gin-wrapper/controllers/admins"
	us "github.com/hiromaily/go-gin-wrapper/controllers/api/users"
	"github.com/hiromaily/go-gin-wrapper/controllers/bases"
	"github.com/hiromaily/go-gin-wrapper/controllers/errors"
	"github.com/hiromaily/go-gin-wrapper/controllers/news"
	//ba "github.com/hiromaily/go-gin-wrapper/libs/basicauth"
	"net/http"
)

//RefererUrl key->request url, value->refer url
var RefererUrls = map[string]string{
	"/login": "login",
	"/user":  "user",
}

//For HTTP
func SetHTTPUrls(r *gin.Engine) {
	/******************************************************************************/
	/******** Return HTML *********************************************************/
	/******************************************************************************/
	//TODO:When http request method is POST, check referer in advance automatically.
	//-----------------------
	//Base(Top Level)
	//-----------------------
	r.GET("/", bases.IndexAction)
	//Redirect
	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	//Login
	r.GET("/login", bases.LoginGetAction)
	r.POST("/login", CheckHttpRefererAndCSRF(), bases.LoginPostAction)

	//Loout
	r.PUT("/logout", bases.LogoutPutAction)   //For Ajax
	r.POST("/logout", bases.LogoutPostAction) //HTML

	//-----------------------
	//News
	//-----------------------
	newsG := r.Group("/news")
	{
		//Top
		newsG.GET("/", news.NewsGetAction)
	}
	//r.GET("/news/", news.NewsGetAction)

	//-----------------------
	//Account(MyPage)
	//-----------------------
	//After login
	accountsG := r.Group("/accounts")
	{
		//Top
		accountsG.GET("/", accounts.AccountsGetAction)
	}
	//r.GET("/accounts/", accounts.AccountsGetAction)

	//-----------------------
	//Admin [BasicAuth()]
	//-----------------------
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"web": "test",
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
	//User
	//-----------------------
	//TODO: which is better to use CheckHttpHeader() for ajax request to check http header.
	// if it's used as middle ware like as below.
	//  r.Use(routes.CheckHttpHeader())
	//  it let us faster to develop instead of a bit less performance.
	users := r.Group("/api/users", CheckHttpHeader())
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
