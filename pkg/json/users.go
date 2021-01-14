package json

import (
	models "github.com/hiromaily/go-gin-wrapper/pkg/model/mysql"
)

// UserListResponse is for user list
type UserListResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Users   []models.UsersSL `json:"users"`
}

// UserIDsResponse is for user IDs
type UserIDsResponse struct {
	Code int   `json:"code"`
	IDs  []int `json:"ids"`
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

// CreateUserIDsJSON is for response for user IDs
func CreateUserIDsJSON(ids []int) *UserIDsResponse {
	userIDs := UserIDsResponse{
		Code: 0,
		IDs:  ids,
	}
	return &userIDs
}

// CreateUserJSON is response for single user
func CreateUserJSON(id int) *UserResponse {
	var (
		code int
		msg  string
	)

	user := UserResponse{
		Code:    code,
		Message: msg,
		ID:      id,
	}
	return &user
}
