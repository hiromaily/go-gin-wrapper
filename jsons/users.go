package jsons

import (
	//u "github.com/hiromaily/golibs/utils"
	models "github.com/hiromaily/go-gin-wrapper/models/mysql"
)

// UserListResponse is for user list
type UserListResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Users   []models.UsersSL `json:"users"`
}

// UserResponse is for single user response
type UserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ID      int    `json:"id"`
}

// CreateUserListJSON is for response for user list
func CreateUserListJSON(users []models.UsersSL) *UserListResponse {
	var (
		code int
		msg  string
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

// CreateUserJSON is response for single user
func CreateUserJSON(id int64) *UserResponse {
	var (
		code int
		msg  string
	)

	user := UserResponse{
		Code:    code,
		Message: msg,
		ID:      int(id),
	}
	return &user
}
