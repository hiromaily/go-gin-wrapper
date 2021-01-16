package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	js "github.com/hiromaily/go-gin-wrapper/pkg/json"
	"github.com/hiromaily/go-gin-wrapper/pkg/model/user"
	jsonresp "github.com/hiromaily/go-gin-wrapper/pkg/server/response/json"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/validator"
	str "github.com/hiromaily/go-gin-wrapper/pkg/strings"
)

// UserRequest is expected request form from user
type UserRequest struct {
	FirstName string `valid:"nonempty,min=3,max=20" field:"first_name" dispName:"First Name" json:"firstName" form:"firstName"`
	LastName  string `valid:"nonempty,min=3,max=20" field:"last_name" dispName:"Last Name" json:"lastName" form:"lastName"`
	Email     string `valid:"nonempty,min=5,max=60" field:"email" dispName:"E-mail" json:"email" form:"email"`
	Password  string `valid:"nonempty,min=8,max=16" field:"password" dispName:"Password" json:"password" form:"password"`
}

// get user parameter and check validation
func (ctl *Controller) getUserParamAndValid(c *gin.Context, data *UserRequest) (err error) {
	// Check token(before send message, pass token)

	contentType := c.Request.Header.Get("Content-Type")
	ctl.logger.Debug("getUserParamAndValid", zap.String("Content-Type", contentType))

	if contentType == "application/json" {
		err = c.BindJSON(data)
	} else {
		// application/x-www-form-urlencoded
		err = c.Bind(data)
	}
	if err != nil {
		return err
	}

	// Validation
	mRet := validator.CheckValidation(data, false)
	if len(mRet) != 0 {
		return errors.New("validation error")
	}

	return nil
}

func (ctl *Controller) getUserParamAndValidForPut(c *gin.Context, data *UserRequest) error {
	//[POST x application/x-www-form-urlencoded] OK
	//[PUT x application/x-www-form-urlencoded] OK
	//[PUT x application/json] NG

	// Param id
	if c.Param("id") == "" {
		return errors.New("missing id on request parameter")
	}

	contentType := c.Request.Header.Get("Content-Type")
	ctl.logger.Debug("getUserParamAndValidForPut", zap.String("Content-Type", contentType))

	var err error
	if contentType == "application/json" {
		err = c.BindJSON(data)
	} else {
		// application/x-www-form-urlencoded
		err = c.Bind(data)
	}
	if err != nil {
		return err
	}
	ctl.logger.Debug("getUserParamAndValidForPut", zap.Any("response body", data))

	// Validation
	if data.FirstName == "" && data.LastName == "" && data.Email == "" && data.Password == "" {
		return errors.New("validation no data error")
	}

	mRet := validator.CheckValidation(data, true)
	if len(mRet) != 0 {
		return errors.New("validation error")
	}

	return nil
}

// insert user
func (ctl *Controller) insertUser(data *UserRequest) (int, error) {
	item := &user.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	}

	// Insert
	return ctl.userRepo.InsertUser(item)
}

// update user
func (ctl *Controller) updateUser(data *UserRequest, id string) (int64, error) {
	item := &user.User{}
	if data.FirstName != "" {
		item.FirstName = data.FirstName
	}
	if data.LastName != "" {
		item.LastName = data.LastName
	}
	if data.Email != "" {
		item.Email = data.Email
	}
	if data.Password != "" {
		item.Password = data.Password
	}

	// Update
	return ctl.userRepo.UpdateUser(item, id)
}

// APIUserListGetAction is get user list [GET]
func (ctl *Controller) APIUserListGetAction(c *gin.Context) {
	ctl.logger.Info("[GET] UsersListGetAction")

	users, err := ctl.userRepo.GetUsers("")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Make json for response and return
	jsonresp.ResponseUserJSON(c, ctl.logger, ctl.cors, 0, js.CreateUserListJSON(users))
}

// ListOptionsAction is preflight request of CORS before get request
//func ListOptionsAction(c *gin.Context) {
//	ctl.logger.Info("[OPTIONS] ListOptionsAction")
//	//TODO: return void??
//	cors.SetHeader(c)
//}

// APIUserInsertPostAction is register for new user [POST]
func (ctl *Controller) APIUserInsertPostAction(c *gin.Context) {
	ctl.logger.Debug("[POST] UserPostAction")

	// Param & Check valid
	var uData UserRequest
	err := ctl.getUserParamAndValid(c, &uData)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// Insert
	id, err := ctl.insertUser(&uData)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	jsonresp.ResponseUserJSON(c, ctl.logger, ctl.cors, 0, js.CreateUserJSON(int(id)))
}

// APIUserGetAction is get specific user [GET]
func (ctl *Controller) APIUserGetAction(c *gin.Context) {
	ctl.logger.Info("[GET] UserGetAction")

	// Param
	// FirstName := c.Query("firstName")
	// ctl.logger.Debug("firstName:", FirstName)
	userID := c.Param("id")
	if userID == "" {
		c.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}

	users, err := ctl.userRepo.GetUsers(userID)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Make json for response and return
	jsonresp.ResponseUserJSON(c, ctl.logger, ctl.cors, 0, js.CreateUserListJSON(users))
}

// APIUserPutAction is update specific user [PUT]
func (ctl *Controller) APIUserPutAction(c *gin.Context) {
	ctl.logger.Info("[PUT] UserPutAction")

	// Param & Check valid
	var uData UserRequest
	err := ctl.getUserParamAndValidForPut(c, &uData)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// Update
	affected, err := ctl.updateUser(&uData, c.Param("id"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if affected == 0 {
		ctl.logger.Debug("there was no updated data.")
	}

	jsonresp.ResponseUserJSON(c, ctl.logger, ctl.cors, 0, js.CreateUserJSON(str.Atoi(c.Param("id"))))
}

// APIUserDeleteAction is delete specific user [DELETE] (work in progress)
func (ctl *Controller) APIUserDeleteAction(c *gin.Context) {
	ctl.logger.Info("[DELETE] UserDeleteAction")
	// check token

	// Param
	if c.Param("id") == "" {
		c.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}

	// Delete
	affected, err := ctl.userRepo.DeleteUser(c.Param("id"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if affected == 0 {
		ctl.logger.Debug("there was no updated data.")
	}

	jsonresp.ResponseUserJSON(c, ctl.logger, ctl.cors, 0, js.CreateUserJSON(str.Atoi(c.Param("id"))))
}

// APIUserIDsGetAction is get user ids [GET]
func (ctl *Controller) APIUserIDsGetAction(c *gin.Context) {
	ctl.logger.Info("[GET] IdsGetAction")

	ids, err := ctl.userRepo.GetUserIDs()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Make json for response and return
	jsonresp.ResponseUserJSON(c, ctl.logger, ctl.cors, 0, js.CreateUserIDsJSON(ids))
}
