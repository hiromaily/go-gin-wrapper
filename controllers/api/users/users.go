package users

import (
	//"encoding/json"
	"errors"
	js "github.com/hiromaily/go-gin-wrapper/jsons"
	jslib "github.com/hiromaily/go-gin-wrapper/libs/response/json"
	models "github.com/hiromaily/go-gin-wrapper/models/mysql"
	lg "github.com/hiromaily/golibs/log"
	tm "github.com/hiromaily/golibs/time"
	u "github.com/hiromaily/golibs/utils"
	"github.com/hiromaily/golibs/validator"
	//gin "gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-gonic/gin"
)

// UserRequest is expected request form from user
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
	return models.GetDB().InsertUser(user)
}

// update user
func updateUser(data *UserRequest, id string) (int64, error) {

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
	//update date
	user.Updated = tm.GetCurrentDateTimeByStr("")

	//Update
	return models.GetDB().UpdateUser(user, id)
}

// ListGetAction is get user list [GET]
func ListGetAction(c *gin.Context) {
	lg.Info("[GET] UsersListGetAction")

	var users []models.UsersSL

	_, err := models.GetDB().GetUserList(&users, "")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	//Make json for response and return
	jslib.RtnUserJSON(c, 0, js.CreateUserListJSON(users))
}

// ListOptionsAction is preflight request of CORS before get request
//func ListOptionsAction(c *gin.Context) {
//	lg.Info("[OPTIONS] ListOptionsAction")
//	//TODO: return void??
//	cors.SetHeader(c)
//}

// InsertPostAction is register for new user [POST]
func InsertPostAction(c *gin.Context) {
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

	jslib.RtnUserJSON(c, 0, js.CreateUserJSON(int(id)))
	return
}

// GetAction is get specific user [GET]
func GetAction(c *gin.Context) {
	lg.Info("[GET] UserGetAction")

	//Param
	//FirstName := c.Query("firstName")
	//lg.Debug("firstName:", FirstName)
	userID := c.Param("id")
	if userID == "" {
		c.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}

	var user models.UsersSL
	b, err := models.GetDB().GetUserList(&user, userID)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	//Make json for response and return
	if b {
		jslib.RtnUserJSON(c, 0, js.CreateUserListJSON([]models.UsersSL{user}))
	} else {
		jslib.RtnUserJSON(c, 0, js.CreateUserListJSON(nil))
	}
}

// PutAction is update specific user [PUT]
func PutAction(c *gin.Context) {
	lg.Info("[PUT] UserPutAction")

	//Param & Check valid
	var uData UserRequest
	err := getUserParamAndValidForPut(c, &uData)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// Update
	affected, err := updateUser(&uData, c.Param("id"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if affected == 0 {
		lg.Debug("there was no updated data.")
	}

	jslib.RtnUserJSON(c, 0, js.CreateUserJSON(u.Atoi(c.Param("id"))))
	return
}

// DeleteAction is delete specific user [DELETE] (work in progress)
func DeleteAction(c *gin.Context) {
	lg.Info("[DELETE] UserDeleteAction")
	//check token

	//Param
	if c.Param("id") == "" {
		c.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}

	//Delete
	affected, err := models.GetDB().DeleteUser(c.Param("id"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if affected == 0 {
		lg.Debug("there was no updated data.")
	}

	jslib.RtnUserJSON(c, 0, js.CreateUserJSON(u.Atoi(c.Param("id"))))
	return
}

// IdsGetAction is get user ids [GET]
func IdsGetAction(c *gin.Context) {
	lg.Info("[GET] IdsGetAction")

	var ids []models.UsersIDs

	err := models.GetDB().GetUserIds(&ids)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	//convert
	//newIds := make([]int, len(ids))
	newIds := []int{}
	for _, id := range ids {
		newIds = append(newIds, id.ID)
	}
	//lg.Debugf("ids: %v", ids)
	//lg.Debugf("newIds: %v", newIds)

	//Make json for response and return
	jslib.RtnUserJSON(c, 0, js.CreateUserIDsJSON(newIds))
}
