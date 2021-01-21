package user

import "time"

// User is for t_user table structure. This is used on Insert, Update method.
type User struct {
	ID        int        `boil:"id" json:"id,omitempty"`
	FirstName string     `boil:"first_name" json:"firstName,omitempty"`
	LastName  string     `boil:"last_name" json:"lastName,omitempty"`
	Email     string     `boil:"email" json:"email,omitempty"`
	Password  string     `boil:"password" json:"password,omitempty"`
	OAuth2    uint8      `boil:"oauth2_type" json:"oauth2_type,omitempty"`
	Created   *time.Time `boil:"updated_at" json:"created_at,omitempty"`
	Updated   *time.Time `boil:"updated_at" json:"updated_at,omitempty"`
}

// UserAuth is response when OAuth2 login is used.
type UserAuth struct {
	ID   int   `boil:"id"`
	Auth uint8 `boil:"oauth2_type"`
}
