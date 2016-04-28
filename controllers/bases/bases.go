package bases

import (
	//"fmt"
	"fmt"
	"github.com/gin-gonic/gin"
	lg "github.com/hiromaily/golibs/log"
	//"github.com/hiromaily/golibs/mysql"
	"github.com/hiromaily/web/libs/session"
	valid "gopkg.in/validator.v2"
	"net/http"
)

type LoginRequest struct {
	//email string `validate:"min=3,max=80,regexp=^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$"`
	Email string `validate:"min=3,max=80"`
	Pass  string `validate:"min=4,max=16"`
}

//Index
func IndexAction(c *gin.Context) {
	//Param

	//Logic

	//View
	c.HTML(http.StatusOK, "bases/index.tmpl", gin.H{
		"title": "Main website",
	})

}

//Login [GET]
func LoginGetAction(c *gin.Context) {
	//fmt.Printf("%#v\n", c)
	//fmt.Printf("method: %s\n", c.Request.Method)
	//fmt.Printf("header: %#v\n", c.Request.Header)
	//fmt.Printf("body: %#v\n", c.Request.Body)

	//Login Logic
	//switch c.Request.Method{
	//case "GET":
	//case "POST":
	//}

	//View
	c.HTML(http.StatusOK, "bases/login.tmpl", gin.H{
		"message": "nothing special",
	})

}

//Login [POST]
func LoginPostAction(c *gin.Context) {
	//TODO:POSTかGETかをどうやって取得する?
	fmt.Printf("%#v\n\n", c)
	fmt.Printf("method: %s\n\n", c.Request.Method)
	fmt.Printf("header: %#v\n\n", c.Request.Header)
	fmt.Printf("body: %#v\n\n", c.Request.Body)

	//Login Logic
	//POSTパラメータを取得
	inputEmail := c.PostForm("inputEmail")
	inputPassword := c.PostForm("inputPassword")
	lg.Debugf("inputEmail : %s\n", inputEmail)
	lg.Debugf("inputPassword : %s\n", inputPassword)
	posted := LoginRequest{Email: inputEmail, Pass: inputPassword}

	//validate and escape
	if errs := valid.Validate(posted); errs != nil {
		//validator.ErrorMap{"Pass":validator.ErrorArray{validator.TextErr{Err:(*errors.errorString)(0xc82012b850)}}}
		lg.Errorf("validation error : %#v\n", errs)

		//エラー関連のメッセージと共にloginページを再表示
		//View
		c.HTML(http.StatusOK, "bases/login.tmpl", gin.H{
			"message": "something error happend",
		})
		return
	}

	//database
	/*
		db := mysql.GetDBInstance()
		sql := "SELECT * FROM t_users WHERE delete_flg=?"
		data, _, err := db.SelectSQLAllField(sql, 0)
		if err != nil {
			panic(err.Error())
		}
	*/

	//Session
	session.SessionStart(c)

	//View(next page)
	c.HTML(http.StatusOK, "bases/login.tmpl", gin.H{
		"message": "next pagee",
	})

}

//func loginPostValid(){
//
//}
