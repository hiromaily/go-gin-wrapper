package json

import (
	"github.com/hiromaily/go-gin-wrapper/pkg/model/user"
)

// UserResponse is response of single user
type UserResponse struct {
	ID int `json:"id"`
}

// UserIDsResponse is response of user IDs
type UserIDsResponse struct {
	IDs []int `json:"ids"`
}

// UserListResponse is response of user list
type UserListResponse struct {
	Users []*user.User `json:"users"`
}

// CreateUserListJSON is for response for user list
func CreateUserListJSON(users []*user.User) *UserListResponse {
	return &UserListResponse{
		Users: users,
	}
}

// CreateUserIDsJSON is for response for user IDs
func CreateUserIDsJSON(ids []int) *UserIDsResponse {
	return &UserIDsResponse{
		IDs: ids,
	}
}

// CreateUserJSON is response for single user
func CreateUserJSON(id int) *UserResponse {
	return &UserResponse{
		ID: id,
	}
}
