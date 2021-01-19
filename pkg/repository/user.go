package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
	"github.com/hiromaily/go-gin-wrapper/pkg/model/user"
	models "github.com/hiromaily/go-gin-wrapper/pkg/models/rdb"
)

// UserRepositorier is UserRepository interface
type UserRepositorier interface {
	IsUserEmail(email, password string) (int, error)
	OAuth2Login(email string) (*user.UserAuth, error)
	GetUserIDs() ([]int, error)
	GetUsers(id string) ([]*user.User, error)
	InsertUser(user *user.User) (int, error)
	UpdateUser(users *user.User, id string) (int64, error)
	DeleteUser(id string) (int64, error)
}

// userRepository is repository for t_users table
type userRepository struct {
	dbConn    *sql.DB
	tableName string
	logger    *zap.Logger
	hash      encryption.Hasher
}

// NewUserRepository returns UserRepositorier
func NewUserRepository(dbConn *sql.DB, logger *zap.Logger, hash encryption.Hasher) UserRepositorier {
	return &userRepository{
		dbConn:    dbConn,
		tableName: "t_user",
		logger:    logger,
		hash:      hash,
	}
}

// IsUserEmail is for User authorization when login
func (u *userRepository) IsUserEmail(email, password string) (int, error) {
	type LoginUser struct {
		ID       int    `boil:"id"`
		Email    string `boil:"email"`
		Password string `boil:"password"`
	}

	ctx := context.Background()

	var user LoginUser
	// sql := "SELECT id, email, password FROM t_users WHERE email=? AND delete_flg=? LIMIT 1"
	err := models.TUsers(
		qm.Select("id, email, password"),
		qm.Where("email=?", email),
		qm.And("delete_flg=?", 0),
	).Bind(ctx, u.dbConn, &user)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call models.TUsers().Bind()")
	}

	// check
	if user.Password != u.hash.Hash(password) {
		return 0, errors.Errorf("password is invalid")
	}
	return user.ID, nil
}

// OAuth2Login is for OAuth2 login
func (u *userRepository) OAuth2Login(email string) (*user.UserAuth, error) {
	// 0:no user -> register and login
	// 1:existing user (google) -> login
	// 2:existing user (no auth or another auth) -> err

	ctx := context.Background()

	var user user.UserAuth
	// sql := "SELECT id, oauth2_flg FROM t_users WHERE email=? AND delete_flg=?"
	err := models.TUsers(
		qm.Select("id, oauth2_flg"),
		qm.Where("email=?", email),
		qm.And("delete_flg=?", 0),
	).Bind(ctx, u.dbConn, &user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.TUsers().Bind()")
	}

	return &user, nil
}

// GetUserIDs is to get user IDs
func (u *userRepository) GetUserIDs() ([]int, error) {
	ctx := context.Background()

	type Response struct {
		ID int `boil:"id"`
	}
	var response []*Response
	// sql := "SELECT id FROM t_users WHERE delete_flg=?"
	err := models.TUsers(
		qm.Select("id"),
		qm.Where("delete_flg=?", 0),
	).Bind(ctx, u.dbConn, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.TUsers().Bind()")
	}

	// convert
	ids := make([]int, 0, len(response))
	for _, val := range response {
		ids = append(ids, val.ID)
	}
	return ids, nil
}

// GetUsers is to get user list
//  parameter id accepts blank
func (u *userRepository) GetUsers(id string) ([]*user.User, error) {
	// sql := "SELECT %s FROM t_user WHERE delete_flg=?"
	ctx := context.Background()

	q := []qm.QueryMod{
		qm.Where("delete_flg=?", 0),
	}
	if id != "" {
		q = append(q, qm.And("id=?", id))
	}
	items, err := models.TUsers(q...).All(ctx, u.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.TUsers().All()")
	}
	// convert
	converted := make([]*user.User, len(items))
	for i, item := range items {
		converted[i] = &user.User{
			ID:        item.ID,
			FirstName: item.FirstName,
			LastName:  item.LastName,
			Email:     item.Email,
			Password:  item.Password,
			OAuth2Flg: item.Oauth2FLG.String,
			Created:   &item.CreatedAt.Time,
			Updated:   &item.UpdatedAt.Time,
		}
	}

	return converted, nil
}

func (u *userRepository) getUserByEmail(email string) (*models.TUser, error) {
	ctx := context.Background()

	q := []qm.QueryMod{
		qm.Where("delete_flg=?", 0),
		qm.And("email=?", email),
	}
	item, err := models.TUsers(q...).One(ctx, u.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.TUsers().All()")
	}
	return item, nil
}

// InsertUser is to insert user
func (u *userRepository) InsertUser(user *user.User) (int, error) {
	item := &models.TUser{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  u.hash.Hash(user.Password),
		Oauth2FLG: null.StringFrom(user.OAuth2Flg),
	}

	ctx := context.Background()
	// sql := "INSERT INTO t_user (first_name, last_name, email, password) VALUES (?,?,?,?)"
	if err := item.Insert(ctx, u.dbConn, boil.Infer()); err != nil {
		return 0, errors.Wrap(err, "failed to call item.Insert()")
	}

	// get Inserted user's ID by email
	use, err := u.getUserByEmail(user.Email)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call getUserByEmail()")
	}
	return use.ID, nil
	//if users.OAuth2Flg != "" {
	//	sql := "INSERT INTO t_users (first_name, last_name, email, password, oauth2_flg) VALUES (?,?,?,?,?)"
	//	// hash password
	//	return us.DB.Insert(sql, users.FirstName, users.LastName, users.Email, encryption.GetMD5Plus(users.Password, ""), users.OAuth2Flg)
	//}
	//// hash password
	//return us.DB.Insert(sql, users.FirstName, users.LastName, users.Email, encryption.GetMD5Plus(users.Password, ""))
}

// UpdateUser is to update user
func (u *userRepository) UpdateUser(users *user.User, id string) (int64, error) {
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{}
	if users.FirstName != "" {
		updCols[models.TUserColumns.FirstName] = users.FirstName
	}
	if users.LastName != "" {
		updCols[models.TUserColumns.LastName] = users.LastName
	}
	if users.Email != "" {
		updCols[models.TUserColumns.Email] = users.Email
	}
	if users.Password != "" {
		updCols[models.TUserColumns.Password] = u.hash.Hash(users.Password)
	}
	updCols[models.TUserColumns.UpdatedAt] = null.TimeFrom(time.Now())

	return models.TUsers(
		qm.Where("id=?", id),
	).UpdateAll(ctx, u.dbConn, updCols)
}

// DeleteUser is to delete user
func (u *userRepository) DeleteUser(id string) (int64, error) {
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{}
	updCols[models.TUserColumns.DeleteFLG] = null.StringFrom("1")
	updCols[models.TUserColumns.UpdatedAt] = null.TimeFrom(time.Now())

	return models.TUsers(
		qm.Where("id=?", id),
	).UpdateAll(ctx, u.dbConn, updCols)
}
