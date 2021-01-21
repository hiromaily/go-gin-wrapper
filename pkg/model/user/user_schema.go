package user

import "time"

// User is customized user response from t_user table
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

// UserAuth is response of OAuth2 login
type UserAuth struct {
	ID   int   `boil:"id"`
	Auth uint8 `boil:"oauth2_type"`
}
