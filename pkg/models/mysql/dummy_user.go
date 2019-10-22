package mysql

// DummyMySQL is dummy object
type DummyMySQL struct{}

// IsUserEmail is for User authorization when trying login
func (d *DummyMySQL) IsUserEmail(email string, password string) (int, error) {
	return 0, nil
}

// OAuth2Login is for OAuth2 login
func (d *DummyMySQL) OAuth2Login(email string) (*UserAuth, error) {
	return nil, nil
}

// GetUserIds is to get user IDs
func (d *DummyMySQL) GetUserIds(users interface{}) error {
	return nil
}

// GetUserList is to get user list
func (d *DummyMySQL) GetUserList(users interface{}, id string) (bool, error) {
	return false, nil
}

// InsertUser is to insert user
func (d *DummyMySQL) InsertUser(users *Users) (int64, error) {
	return 0, nil
}

// UpdateUser is to update user
func (d *DummyMySQL) UpdateUser(users *Users, id string) (int64, error) {
	return 0, nil
}

// DeleteUser is to delete user
func (d *DummyMySQL) DeleteUser(id string) (int64, error) {
	return 0, nil
}
