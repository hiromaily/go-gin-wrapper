package controllers

import (
	//"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	js "github.com/hiromaily/go-gin-wrapper/pkg/json"
	models "github.com/hiromaily/go-gin-wrapper/pkg/model/mysql"
	jslib "github.com/hiromaily/go-gin-wrapper/pkg/server/response/json"
	lg "github.com/hiromaily/golibs/log"
	tm "github.com/hiromaily/golibs/time"
	u "github.com/hiromaily/golibs/utils"
	"github.com/hiromaily/golibs/validator"
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
	lg.Debug("Content-Type is ", contentType)

	if contentType == "application/json" {
		err = c.BindJSON(data)
	} else {
		// application/x-www-form-urlencoded
		err = c.Bind(data)
	}
	if err != nil {
		return err
	}

	lg.Debug("Request Body: %#v", data)

	// Validation
	mRet := validator.CheckValidation(data, false)
	lg.Debug(mRet)
	if len(mRet) != 0 {
		return errors.New("validation error")
	}

	return nil
}

func (ctl *Controller) getUserParamAndValidForPut(c *gin.Context, data *UserRequest) (err error) {
	// Param id
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
		// application/x-www-form-urlencoded
		err = c.Bind(data)
	}
	if err != nil {
		return err
	}

	lg.Debug("Request Body: %#v", data)

	// Validation
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
func (ctl *Controller) insertUser(data *UserRequest) (int64, error) {
	user := &models.Users{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
	}

	// Insert
	return ctl.db.InsertUser(user)
}

// update user
func (ctl *Controller) updateUser(data *UserRequest, id string) (int64, error) {
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
	// update date
	user.Updated = tm.GetCurrentDateTimeByStr("")

	// Update
	return ctl.db.UpdateUser(user, id)
}

// APIUserListGetAction is get user list [GET]
func (ctl *Controller) APIUserListGetAction(c *gin.Context) {
	lg.Info("[GET] UsersListGetAction")

	var users []models.UsersSL

	_, err := ctl.db.GetUserList(&users, "")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Make json for response and return
	jslib.ResponseUserJSON(c, ctl.cors, 0, js.CreateUserListJSON(users))
}

// ListOptionsAction is preflight request of CORS before get request
//func ListOptionsAction(c *gin.Context) {
//	lg.Info("[OPTIONS] ListOptionsAction")
//	//TODO: return void??
//	cors.SetHeader(c)
//}

// APIUserInsertPostAction is register for new user [POST]
func (ctl *Controller) APIUserInsertPostAction(c *gin.Context) {
	lg.Debug("[POST] UserPostAction")

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

	jslib.ResponseUserJSON(c, ctl.cors, 0, js.CreateUserJSON(int(id)))
}

// APIUserGetAction is get specific user [GET]
func (ctl *Controller) APIUserGetAction(c *gin.Context) {
	lg.Info("[GET] UserGetAction")

	// Param
	// FirstName := c.Query("firstName")
	// lg.Debug("firstName:", FirstName)
	userID := c.Param("id")
	if userID == "" {
		c.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}

	var user models.UsersSL
	b, err := ctl.db.GetUserList(&user, userID)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Make json for response and return
	if b {
		jslib.ResponseUserJSON(c, ctl.cors, 0, js.CreateUserListJSON([]models.UsersSL{user}))
	} else {
		jslib.ResponseUserJSON(c, ctl.cors, 0, js.CreateUserListJSON(nil))
	}
}

// APIUserPutAction is update specific user [PUT]
func (ctl *Controller) APIUserPutAction(c *gin.Context) {
	lg.Info("[PUT] UserPutAction")

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
		lg.Debug("there was no updated data.")
	}

	jslib.ResponseUserJSON(c, ctl.cors, 0, js.CreateUserJSON(u.Atoi(c.Param("id"))))
}

// APIUserDeleteAction is delete specific user [DELETE] (work in progress)
func (ctl *Controller) APIUserDeleteAction(c *gin.Context) {
	lg.Info("[DELETE] UserDeleteAction")
	// check token

	// Param
	if c.Param("id") == "" {
		c.AbortWithError(400, errors.New("missing id on request parameter"))
		return
	}

	// Delete
	affected, err := ctl.db.DeleteUser(c.Param("id"))
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	if affected == 0 {
		lg.Debug("there was no updated data.")
	}

	jslib.ResponseUserJSON(c, ctl.cors, 0, js.CreateUserJSON(u.Atoi(c.Param("id"))))
}

// APIUserIDsGetAction is get user ids [GET]
func (ctl *Controller) APIUserIDsGetAction(c *gin.Context) {
	lg.Info("[GET] IdsGetAction")

	var ids []models.UsersIDs
	err := ctl.db.GetUserIds(&ids)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// convert
	// newIds := make([]int, len(ids))
	newIds := []int{}
	for _, id := range ids {
		newIds = append(newIds, id.ID)
	}
	// lg.Debugf("ids: %v", ids)
	// lg.Debugf("newIds: %v", newIds)

	// Make json for response and return
	jslib.ResponseUserJSON(c, ctl.cors, 0, js.CreateUserIDsJSON(newIds))
}
