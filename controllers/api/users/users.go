package users

import (
	//"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	js "github.com/hiromaily/go-gin-wrapper/jsons"
	jslib "github.com/hiromaily/go-gin-wrapper/libs/json"
	"github.com/hiromaily/go-gin-wrapper/models"
	lg "github.com/hiromaily/golibs/log"
	"github.com/hiromaily/golibs/validator"
)

//{'firstName':'kentaro','lastName':'asakura','email':'cccc@aa.jp', 'password':'testtest'};
type UserRequest struct {
	FirstName string `valid:"nonempty,min=3,max=20" field:"first_name" dispName:"First Name" json:"firstName" form:"firstName"`
	LastName  string `valid:"nonempty,min=3,max=20" field:"last_name" dispName:"Last Name" json:"lastName" form:"lastName"`
	Email     string `valid:"nonempty,min=5,max=60" field:"email" dispName:"E-mail" json:"email" form:"email"`
	Password  string `valid:"nonempty,min=8,max=16" field:"password" dispName:"Password" json:"password" form:"password"`
}

// get user parameter and check validation
func getUserParamAndValid(c *gin.Context, data *UserRequest) (err error) {
	//Check token(before send message, pass token)

	contentType := c.Request.Header.Get("Content-Type")
	lg.Debug("Content-Type is ", contentType)

	if contentType == "application/json" {
		err = c.BindJSON(data)
	} else {
		//application/x-www-form-urlencoded
		err = c.Bind(data)
	}
	if err != nil {
		return err
	}

	lg.Debug("Request Body: %#v", data)

	//Validation
	mRet := validator.CheckValidation(data, false)
	lg.Debug(mRet)
	if len(mRet) != 0 {
		return errors.New("validation error")
	}

	return nil
}

func getUserParamAndValidForPut(c *gin.Context, data *UserRequest) (err error) {
	//Param id
	if c.Param("id") == "" {
		return errors.New("missing id on request parameter")
	}

	//[POST x application/x-www-form-urlencoded] OK
	//[PUT x application/x-www-form-urlencoded] OK
	//[PUT x application/json] NG
	//lg.Debug("firstName:", c.PostForm("firstName"))
	//lg.Debug("email:", c.PostForm("email"))

	contentType := c.Request.Header.Get("Content-Type")
	lg.Debug("Content-Type is ", contentType)

	if contentType == "application/json" {
		err = c.BindJSON(data)
	} else {
		//application/x-www-form-urlencoded
		err = c.Bind(data)
	}
	if err != nil {
		return err
	}

	lg.Debug("Request Body: %#v", data)

	//Validation
	if data.FirstName == "" && data.LastName == "" && data.Email == "" && data.Password == "" {
		return errors.New("validation no data error")
	}

	mRet := validator.CheckValidation(data, true)
	lg.Debug(mRet)
	if len(mRet) != 0 {
		return errors.New("validation error")
	}

	return nil
}

// insert user
func insertUser(c *gin.Context, data *UserRequest) (int64, error) {

	user := &models.Users{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	}

	//Insert
	return models.GetDBInstance().InsertUser(user)
}

// update user
func updateUser(c *gin.Context, data *UserRequest, id string) (int64, error) {

	user := &models.Users{}
	if data.FirstName != "" {
		user.FirstName = data.FirstName
	}
	if data.LastName != "" {
		user.LastName = data.LastName
	}
	if data.Email != "" {
		user.Email = data.Email
	}
	if data.Password != "" {
		user.Password = data.Password
	}

	//Update
	return models.GetDBInstance().UpdateUser(user, id)
}

//Users: get list [GET]
func UsersListGetAction(c *gin.Context) {
	lg.Debug("[GET] UsersListGetAction")

	//Param
	//FirstName := c.Query("firstName")
	//lg.Debug("firstName:", FirstName)

	mapRet, err := models.GetDBInstance().GetUserList("")
	if err != nil {
		c.AbortWithError(500, err)
		return
	} else {
		//Make json for response and return
		jslib.RtnUserJson(c, 0, js.CreateUserListJson(mapRet))
		return
	}
}

//Users: register for new user [POST]
func UserPostAction(c *gin.Context) {
	lg.Debug("[POST] UserPostAction")

	//Param & Check valid
	var uData UserRequest
	err := getUserParamAndValid(c, &uData)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// Insert
	id, err := insertUser(c, &uData)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	jslib.RtnUserJson(c, 0, js.CreateUserJson(id))
	return
}

//Users: get specific user [GET]
func UserGetAction(c *gin.Context) {
	lg.Debug("[GET] UserGetAction")

	//Param
	//FirstName := c.Query("firstName")
	//lg.Debug("firstName:", FirstName)
	userId := c.Param("id")
	if userId == "" {
		c.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}

	mapRet, err := models.GetDBInstance().GetUserList(userId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	} else {
		//Make json for response and return
		jslib.RtnUserJson(c, 0, js.CreateUserListJson(mapRet))
		return
	}
}

//Users: update specific user [PUT]
func UserPutAction(c *gin.Context) {
	lg.Debug("[PUT] UserPutAction")

	//Param & Check valid
	var uData UserRequest
	err := getUserParamAndValidForPut(c, &uData)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// Update
	_, err = updateUser(c, &uData, c.Param("id"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	jslib.RtnUserJson(c, 0, js.CreateUserJson(0))
	return
}

//Users: delete specific user [DELETE] (work in progress)
func UserDeleteAction(c *gin.Context) {
	lg.Debug("[DELETE] UserDeleteAction")
	//check token

	//Param
	if c.Param("id") == "" {
		c.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}

	//Delete
	err := models.GetDBInstance().DeleteUser(c.Param("id"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	jslib.RtnUserJson(c, 0, js.CreateUserJson(0))
	return
}
