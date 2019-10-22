package mysql

// Users is for t_users table structure. This is used on Insert, Update method.
type Users struct {
	ID        int    `column:"user_id" json:"id"`
	FirstName string `column:"first_name" json:"firstName"`
	LastName  string `column:"last_name" json:"lastName"`
	Email     string `column:"email" json:"email"`
	Password  string `column:"password" json:"password"`
	OAuth2Flg string `column:"oauth2_flg" json:"oauth2_flg"`
	//DeleteFlg string    `column:"delete_flg"       db:"delete_flg"`
	//Created   time.Time `column:"create_datetime"  db:"create_datetime"`
	Updated string `column:"update_datetime" json:"update"`
}

// UsersSL is for t_users table structure. This is used on Select method.
type UsersSL struct {
	ID        int    `column:"user_id" json:"id"`
	FirstName string `column:"first_name" json:"firstName"`
	LastName  string `column:"last_name" json:"lastName"`
	Email     string `column:"email" json:"email"`
	Updated   string `column:"update_datetime" json:"updated"`
}

// UsersIDs is for t_user table structure. This is for only user_id.
type UsersIDs struct {
	ID int `column:"user_id" json:"id"`
}

// UserAuth is response when OAuth2 login is used.
type UserAuth struct {
	ID   int
	Auth string
}
