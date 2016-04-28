package routes

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hiromaily/web/controllers/bases"
	"github.com/hiromaily/web/controllers/errors"
	us "github.com/hiromaily/web/controllers/users"
	"net/http"
)

func SetUrls(r *gin.Engine) {

	/* Return HTML */
	//Redirect
	r.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	//Base(Top Level)
	r.GET("/", bases.IndexAction)
	r.GET("/login", bases.LoginGetAction)
	r.POST("/login", bases.LoginPostAction)

	//Account(MyPage)

	//Admin

	/* Error HTML */
	r.NoRoute(errors.Error404Action)
	r.NoMethod(errors.Error405Action)
	//TODO:500Error

	/* REST API */
	//User
	users := r.Group("/api/users")
	{
		//users.Handle("GET","/", users.UsersAction)
		users.GET("/", us.UsersAction)        //一覧取得
		users.POST("/", us.UsersAction)       //新規登録
		users.GET("/:id/", us.UsersAction)    //特定ユーザーの取得
		users.PUT("/:id/", us.UsersAction)    //特定ユーザーの更新
		users.DELETE("/:id/", us.UsersAction) //特定ユーザーの削除
		//必須じゃないパラメータは*XXXXと記述する /user/:name/*action
	}

}
