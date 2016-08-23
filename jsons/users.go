package jsons

import (
	//u "github.com/hiromaily/golibs/utils"
	models "github.com/hiromaily/go-gin-wrapper/models/mysql"
)

/***************** User Json *****************/
// User List
type UserListResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Users   []models.UsersSL `json:"users"`
}

// User
/*
type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
*/

// User Response
type UserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Id      int    `json:"id"`
}

// Response for user list
//func CreateUserListJson(data []map[string]interface{}) *UserListResponse {
func CreateUserListJson(users []models.UsersSL) *UserListResponse {
	var (
		code int    = 0
		msg  string = ""
	)

	if users == nil {
		code = 1
		msg = "no users"
	}

	userList := UserListResponse{
		Code:    code,
		Message: msg,
		Users:   users,
	}
	return &userList
}

// Response for user
func CreateUserJson(id int64) *UserResponse {
	var (
		code int    = 0
		msg  string = ""
	)

	user := UserResponse{
		Code:    code,
		Message: msg,
		Id:      int(id),
	}
	return &user
}
