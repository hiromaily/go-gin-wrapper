package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/model/user"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/ginbinder"
	jsonresp "github.com/hiromaily/go-gin-wrapper/pkg/server/response/json"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/validator"
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

func validateAtoi(val, name string) (int, error) {
	if val == "" {
		return 0, errors.Errorf("%s is missing", name)
	}
	id, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.Errorf("invalid %s: %s", name, val)
	}
	return id, nil
}

func validateUserRequest(ctx *gin.Context, data *UserRequest) error {
	if err := ginbinder.Bind(ctx, data); err != nil {
		return err
	}

	result := validator.Validate(data, false)
	if len(result) != 0 {
		return errors.New("validation error")
	}
	return nil
}

func validateUserRequestUpdate(ctx *gin.Context, data *UserRequest) (int, error) {
	// [POST x application/x-www-form-urlencoded] OK
	// [PUT x application/x-www-form-urlencoded] OK
	// [PUT x application/json] NG
	id, err := validateAtoi(ctx.Param("id"), "id")
	if err != nil {
		return 0, err
	}
	if err = ginbinder.Bind(ctx, data); err != nil {
		return 0, err
	}

	if data.FirstName == "" && data.LastName == "" && data.Email == "" && data.Password == "" {
		return 0, errors.New("validation error")
	}
	result := validator.Validate(data, true)
	if len(result) != 0 {
		return 0, errors.New("validation error")
	}

	return id, nil
}

func (ctl *controller) insertUser(data *UserRequest) (int, error) {
	item := &user.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	}
	return ctl.userRepo.InsertUser(item)
}

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
	return ctl.userRepo.UpdateUser(item, id)
}

// APIUserListGetAction returns user list [GET]
func (ctl *controller) APIUserListGetAction(ctx *gin.Context) {
	ctl.logger.Info("[GET] UsersListGetAction")

	users, err := ctl.userRepo.GetUsers("")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// json response
	ctx.JSON(http.StatusOK, jsonresp.CreateUserListJSON(users))
}

// APIUserInsertPostAction inserts new user [POST]
func (ctl *controller) APIUserInsertPostAction(ctx *gin.Context) {
	ctl.logger.Debug("[POST] UserPostAction")

	var userRequest UserRequest
	if err := validateUserRequest(ctx, &userRequest); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := ctl.insertUser(&userRequest)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// json response
	ctx.JSON(http.StatusOK, jsonresp.CreateUserJSON(id))
}

// APIUserGetAction returns specific user [GET]
func (ctl *controller) APIUserGetAction(ctx *gin.Context) {
	ctl.logger.Info("[GET] UserGetAction")

	id := ctx.Param("id")
	_, err := validateAtoi(id, "id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	users, err := ctl.userRepo.GetUsers(id)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// json response
	ctx.JSON(http.StatusOK, jsonresp.CreateUserListJSON(users))
}

// APIUserPutAction updates specific user [PUT]
func (ctl *controller) APIUserPutAction(ctx *gin.Context) {
	ctl.logger.Info("[PUT] UserPutAction")

	var userRequest UserRequest
	userID, err := validateUserRequestUpdate(ctx, &userRequest)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	affected, err := ctl.updateUser(&userRequest, userID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if affected == 0 {
		ctl.logger.Debug("nothing updated")
	}

	// json response
	ctx.JSON(http.StatusOK, jsonresp.CreateUserJSON(userID))
}

// APIUserDeleteAction deletes specific user [DELETE]
func (ctl *controller) APIUserDeleteAction(ctx *gin.Context) {
	ctl.logger.Info("[DELETE] UserDeleteAction")

	userID, err := validateAtoi(ctx.Param("id"), "id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	affected, err := ctl.userRepo.DeleteUser(userID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if affected == 0 {
		ctl.logger.Debug("nothing deleted")
	}

	// json response
	ctx.JSON(http.StatusOK, jsonresp.CreateUserJSON(userID))
}

// APIUserIDsGetAction returns user ids [GET]
func (ctl *controller) APIUserIDsGetAction(ctx *gin.Context) {
	ctl.logger.Info("[GET] UserIDsGetAction")

	ids, err := ctl.userRepo.GetUserIDs()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// json response
	ctx.JSON(http.StatusOK, jsonresp.CreateUserIDsJSON(ids))
}
