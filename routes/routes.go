package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	//"reflect"
	"time"
)

/* sample code */

// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func loginEndpoint(c *gin.Context) {
	name := c.PostForm("name")
	message := c.PostForm("message")

	fmt.Printf("name: %s; message: %s \n", name, message)
}

func submitEndpoint(c *gin.Context) {
	name := c.PostForm("name")
	message := c.PostForm("message")

	fmt.Printf("name: %s; message: %s \n", name, message)
}

func readEndpoint(c *gin.Context) {
	name := c.PostForm("name")
	message := c.PostForm("message")

	fmt.Printf("name: %s; message: %s \n", name, message)
}

func analyticsEndpoint(c *gin.Context) {
	name := c.PostForm("name")
	message := c.PostForm("message")

	fmt.Printf("name: %s; message: %s \n", name, message)
}

//Loger
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func New(router *gin.Engine) {
	//set router
	//*gin.Engine
	/*
	   router.GET("/someGet", getting)
	   router.POST("/somePost", posting)
	   router.PUT("/somePut", putting)
	   router.DELETE("/someDelete", deleting)
	   router.PATCH("/somePatch", patching)
	   router.HEAD("/someHead", head)
	   router.OPTIONS("/someOptions", options)
	*/

	setGetRouter(router)
	setPostRouter(router)
	setBasicAuthRouter(router)

	//10.Serving static files
	//   HTML rendering

	//11.Redirects

}

//set router for Get Method
func setGetRouter(router *gin.Engine) {
	//func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	//return group.handle("GET", relativePath, handlers)

	//index
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	//1.Parameters in path
	// This handler will match /user/john but will not match neither /user/ or /user
	// e.g. /user/hiroki or /user/jiro
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/join/
	// e.g. /user/hiroki/abcd or /user/hiroki/abcd/efgh
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	//2.Querystring parameters
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	// e.g. /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	//9.XML and JSON rendering
	// gin.H is a shortcut for map[string]interface{}
	// curl -X GET http://localhost:9999/someJSON
	router.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	// curl -X GET http://localhost:9999/moreJSON
	router.GET("/moreJSON", func(c *gin.Context) {
		// You also can use a struct
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// Note that msg.Name becomes "user" in the JSON
		// Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
		c.JSON(http.StatusOK, msg)
	})

	// curl -X GET http://localhost:9999/someXML
	router.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	//12.Custom Middleware
	// curl -X GET http://localhost:9999/logger
	// curl -X GET "http://localhost:9999/logger?example=test"
	router.Use(Logger())
	router.GET("/logger", func(c *gin.Context) {
		example := c.MustGet("example").(string)
		//Returns the value for the given key if it exists, otherwise it panics.

		// it would print: "12345"
		log.Println(example)
	})

}

//set router for Post Method
func setPostRouter(router *gin.Engine) {
	//3.POST data
	//curl -d message=goodnight -d name=taro http://localhost:9999/form_post
	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		name := c.DefaultPostForm("name", "anonymous")

		//JSONをレスポンスとして返す
		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"name":    name,
		})
	})

	//4.query + post form
	//curl -d message=goodnight -d name=taro "http://localhost:9999/postget?id=123&page=1"
	router.POST("/postget", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s \n", id, page, name, message)
	})

	//5.Grouping routes
	// Simple group: v1
	// e.g. /v1/loginってこと？
	//curl -d message=goodnight -d name=taro http://localhost:9999/v1/login
	//curl -d message=goodnight -d name=taro http://localhost:9999/v1/submit
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// Simple group: v2
	/*
		v2 := router.Group("/v2")
		{
			v2.POST("/login", loginEndpoint)
			v2.POST("/submit", submitEndpoint)
			v2.POST("/read", readEndpoint)
		}
	*/

	//6.Model binding and validation
	// Example for binding JSON ({"user": "manu", "password": "123"})
	// JSONをリクエストせねばならないかと。
	//curl -v -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"user": "manu", "password": "123"}' http://localhost:9999/loginJSON
	router.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if c.BindJSON(&json) == nil {
			if json.User == "manu" && json.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})

	//7. Example for binding a HTML form (user=manu&password=123)
	//curl -d user=manu -d password=123 http://localhost:9999/loginForm
	//curl -v --form user=manu --form password=123 http://localhost:9999/loginForm
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if c.Bind(&form) == nil {
			if form.User == "manu" && form.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})

	//8.Multipart/Urlencoded binding
	//curl -v --form user=taro --form password=abcd http://localhost:9999/loginpage
	router.POST("/loginpage", func(c *gin.Context) {
		// you can bind multipart form with explicit binding declaration:
		// c.BindWith(&form, binding.Form)
		// or you can simply use autobinding with Bind method:
		var form Login
		// in this case proper binding will be automatically selected
		if c.Bind(&form) == nil {
			if form.User == "taro" && form.Password == "abcd" {
				c.JSON(200, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(401, gin.H{"status": "unauthorized"})
			}
		}
	})

}

//set router for Post Method
func setBasicAuthRouter(router *gin.Engine) {
	//13.Using BasicAuth() middleware(ベーシック認証)
	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets endpoint
	// hit "localhost:8080/admin/secrets
	// curl --user austin:1234 http://localhost:9999/admin/secrets
	// これ、ベーシック認証じゃね？ブラウザでアクセス
	authorized.GET("/secrets", func(c *gin.Context) {
		// get user, it was setted by the BasicAuth middleware
		//const AuthUserKey = "user"
		//func (c *Context) MustGet(key string) interface{} {
		log.Println(gin.AuthUserKey) //user
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	//14.Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	//curl --user foo:bar -d user=foo -d password=bar -d value=123 http://localhost:9999/admin2
	authorized2 := router.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized2.POST("/admin2", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		//fmt.Printf("%f", user)
		fmt.Println(user)

		// Parse JSON
		var json struct {
			Value string `form:"value"　json:"value" binding:"required"`
		}
		//jsonの場合は、BindJSONを使わないといけないのでは？
		if c.Bind(&json) == nil {
			fmt.Println("test3")
			//DB[user] = json.Value
			c.JSON(200, gin.H{"status": "ok"})
		}
	})

}
