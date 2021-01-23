package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/model/user"
	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
	jsonresp "github.com/hiromaily/go-gin-wrapper/pkg/server/response/json"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/validator"
	str "github.com/hiromaily/go-gin-wrapper/pkg/strings"
)

// APIUser interface
type APIUser interface {
	APIUserListGetAction(ctx *gin.Context)
	APIUserInsertPostAction(ctx *gin.Context)
	APIUserGetAction(ctx *gin.Context)
	APIUserPutAction(ctx *gin.Context)
	APIUserDeleteAction(ctx *gin.Context)
	APIUserIDsGetAction(ctx *gin.Context)
}

// UserRequest is expected request payload
type UserRequest struct {
	FirstName string `valid:"nonempty,min=3,max=20" field:"first_name" dispName:"First Name" json:"firstName" form:"firstName"`
	LastName  string `valid:"nonempty,min=3,max=20" field:"last_name" dispName:"Last Name" json:"lastName" form:"lastName"`
	Email     string `valid:"nonempty,min=5,max=60" field:"email" dispName:"E-mail" json:"email" form:"email"`
	Password  string `valid:"nonempty,min=8,max=16" field:"password" dispName:"Password" json:"password" form:"password"`
}

// get user parameter and check validation
func (ctl *controller) getUserParamAndValid(ctx *gin.Context, data *UserRequest) (err error) {
	// Check token(before send message, pass token)

	// FIXME: middleware has responsibility
	contentType := ctx.Request.Header.Get("Content-Type")
	ctl.logger.Debug("getUserParamAndValid", zap.String("Content-Type", contentType))

	if contentType == "application/json" {
		err = ctx.BindJSON(data)
	} else {
		// application/x-www-form-urlencoded
		err = ctx.Bind(data)
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

func (ctl *controller) getUserParamAndValidForPut(ctx *gin.Context, data *UserRequest) (int, error) {
	//[POST x application/x-www-form-urlencoded] OK
	//[PUT x application/x-www-form-urlencoded] OK
	//[PUT x application/json] NG

	// Param id
	if ctx.Param("id") == "" {
		return 0, errors.New("missing id on request parameter")
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return 0, errors.Errorf("invalid id: %s", ctx.Param("id"))
	}

	contentType := ctx.Request.Header.Get("Content-Type")
	ctl.logger.Debug("getUserParamAndValidForPut", zap.String("Content-Type", contentType))

	if contentType == "application/json" {
		err = ctx.BindJSON(data)
	} else {
		// application/x-www-form-urlencoded
		err = ctx.Bind(data)
	}
	if err != nil {
		return 0, err
	}
	ctl.logger.Debug("getUserParamAndValidForPut", zap.Any("response body", data))

	// Validation
	if data.FirstName == "" && data.LastName == "" && data.Email == "" && data.Password == "" {
		return 0, errors.New("validation no data error")
	}

	mRet := validator.CheckValidation(data, true)
	if len(mRet) != 0 {
		return 0, errors.New("validation error")
	}

	return id, nil
}

// insert user
func (ctl *controller) insertUser(data *UserRequest) (int, error) {
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
func (ctl *controller) updateUser(data *UserRequest, id int) (int64, error) {
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
func (ctl *controller) APIUserListGetAction(ctx *gin.Context) {
	ctl.logger.Info("[GET] UsersListGetAction")

	users, err := ctl.userRepo.GetUsers("")
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	hh.SetResponseHeader(ctx, ctl.logger)
	// FIXME
	//if ctl.corsConf.Enabled && ctx.Request.Method == "GET" {
	//	cors.SetHeader(ctx)
	//}
	// json response
	jsonresp.ResponseUserJSON(ctx, http.StatusOK, jsonresp.CreateUserListJSON(users))
}

// ListOptionsAction is preflight request of CORS before get request
//func ListOptionsAction(ctx *gin.Context) {
//	ctl.logger.Info("[OPTIONS] ListOptionsAction")
//	//TODO: return void??
//	cors.SetHeader(ctx)
//}

// APIUserInsertPostAction is register for new user [POST]
func (ctl *controller) APIUserInsertPostAction(ctx *gin.Context) {
	ctl.logger.Debug("[POST] UserPostAction")

	// Param & Check valid
	var uData UserRequest
	err := ctl.getUserParamAndValid(ctx, &uData)
	if err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	// Insert
	id, err := ctl.insertUser(&uData)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	hh.SetResponseHeader(ctx, ctl.logger)
	// FIXME
	//if ctl.corsConf.Enabled && ctx.Request.Method == "GET" {
	//	cors.SetHeader(ctx)
	//}
	// json response
	jsonresp.ResponseUserJSON(ctx, http.StatusOK, jsonresp.CreateUserJSON(id))
}

// APIUserGetAction is get specific user [GET]
func (ctl *controller) APIUserGetAction(ctx *gin.Context) {
	ctl.logger.Info("[GET] UserGetAction")

	// Param
	// FirstName := ctx.Query("firstName")
	// ctl.logger.Debug("firstName:", FirstName)
	userID := ctx.Param("id")
	if userID == "" {
		ctx.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}

	users, err := ctl.userRepo.GetUsers(userID)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	hh.SetResponseHeader(ctx, ctl.logger)
	// FIXME
	//if ctl.corsConf.Enabled && ctx.Request.Method == "GET" {
	//	cors.SetHeader(ctx)
	//}
	// json response
	jsonresp.ResponseUserJSON(ctx, http.StatusOK, jsonresp.CreateUserListJSON(users))
}

// APIUserPutAction is update specific user [PUT]
func (ctl *controller) APIUserPutAction(ctx *gin.Context) {
	ctl.logger.Info("[PUT] UserPutAction")

	// Param & Check valid
	var uData UserRequest
	id, err := ctl.getUserParamAndValidForPut(ctx, &uData)
	if err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	// Update
	affected, err := ctl.updateUser(&uData, id)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	if affected == 0 {
		ctl.logger.Debug("there was no updated data.")
	}

	hh.SetResponseHeader(ctx, ctl.logger)
	// FIXME
	//if ctl.corsConf.Enabled && ctx.Request.Method == "GET" {
	//	cors.SetHeader(ctx)
	//}
	// json response
	jsonresp.ResponseUserJSON(ctx, http.StatusOK, jsonresp.CreateUserJSON(str.Atoi(ctx.Param("id"))))
}

// APIUserDeleteAction is delete specific user [DELETE] (work in progress)
func (ctl *controller) APIUserDeleteAction(ctx *gin.Context) {
	ctl.logger.Info("[DELETE] UserDeleteAction")
	// check token

	// Param
	if ctx.Param("id") == "" {
		ctx.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(400, errors.Errorf("invalid id: %s", ctx.Param("id")))
		return
	}

	// Delete
	affected, err := ctl.userRepo.DeleteUser(id)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	if affected == 0 {
		ctl.logger.Debug("there was no updated data.")
	}

	hh.SetResponseHeader(ctx, ctl.logger)
	// FIXME
	//if ctl.corsConf.Enabled && ctx.Request.Method == "GET" {
	//	cors.SetHeader(ctx)
	//}
	// json response
	jsonresp.ResponseUserJSON(ctx, http.StatusOK, jsonresp.CreateUserJSON(str.Atoi(ctx.Param("id"))))
}

// APIUserIDsGetAction is get user ids [GET]
func (ctl *controller) APIUserIDsGetAction(ctx *gin.Context) {
	ctl.logger.Info("[GET] IdsGetAction")

	ids, err := ctl.userRepo.GetUserIDs()
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	hh.SetResponseHeader(ctx, ctl.logger)
	// FIXME
	//if ctl.corsConf.Enabled && ctx.Request.Method == "GET" {
	//	cors.SetHeader(ctx)
	//}
	// json response
	jsonresp.ResponseUserJSON(ctx, http.StatusOK, jsonresp.CreateUserIDsJSON(ids))
}
