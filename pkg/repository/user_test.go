// +build integration

package repository

import (
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
	"github.com/hiromaily/go-gin-wrapper/pkg/logger"
	"github.com/hiromaily/go-gin-wrapper/pkg/model/user"
	"github.com/hiromaily/go-gin-wrapper/pkg/storage/mysql"
)

// Note:
// - test data is in ./test/sql/data_gogin-test.sql
// - `settings.toml` is used for this unittest
// - result of md5 hash can be found by command `tool-md5`, see Makefile
// - run the below commands before this unittest
// $ docker-compose mysql
// $ make setup-testdb

var userRepo UserRepositorier

func getUserRepo(t *testing.T) UserRepositorier {
	if userRepo == nil {
		// config
		conf, err := config.GetConf("settings.toml")
		if err != nil {
			t.Fatal(err)
		}
		// db
		dbConn, err := mysql.NewMySQL(conf.MySQL.Test)
		if err != nil {
			t.Fatal(err)
		}
		//
		userRepo = NewUserRepository(
			dbConn,
			logger.NewZapLogger(conf.Logger),
			encryption.NewMD5(conf.Hash.Salt1, conf.Hash.Salt2),
		)
	}
	return userRepo
}

func TestIsUserEmail(t *testing.T) {
	repo := getUserRepo(t)

	type args struct {
		email    string
		password string
	}
	type want struct {
		userID int
		isErr  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				email:    "foobar@gogin.com",
				password: "password",
			},
			want: want{
				userID: 1,
				isErr:  false,
			},
		},
		{
			name: "happy path 2",
			args: args{
				email:    "mark@gogin.com",
				password: "secret-string",
			},
			want: want{
				userID: 2,
				isErr:  false,
			},
		},
		{
			name: "no such email",
			args: args{
				email:    "no-user@gogin.com",
				password: "secret-string",
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "password is wrong",
			args: args{
				email:    "foobar@gogin.com",
				password: "wrong-password",
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "no email",
			args: args{
				email:    "",
				password: "xxxxx",
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "no password",
			args: args{
				email:    "foobar@gogin.com",
				password: "",
			},
			want: want{
				isErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := repo.IsUserEmail(tt.args.email, tt.args.password)
			if (err != nil) != tt.want.isErr {
				t.Errorf("IsUserEmail() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				if err != nil {
					t.Log(err)
				}
				return
			}
			if err != nil {
				return
			}
			if userID != tt.want.userID {
				t.Errorf("IsUserEmail(): userID = %d, want %d", userID, tt.want.userID)
			}
		})
	}
}

func TestOAuth2Login(t *testing.T) {
	repo := getUserRepo(t)

	type args struct {
		email    string
	}
	type want struct {
		userAuth *user.UserAuth
		isErr    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				email: "chk-authtype1@gogin.com",
			},
			want: want{
				userAuth: &user.UserAuth{
					ID:   3,
					Auth: 1,
				},
				isErr: false,
			},
		},
		{
			name: "happy path 2",
			args: args{
				email: "chk-authtype2@gogin.com",
			},
			want: want{
				userAuth: &user.UserAuth{
					ID:   4,
					Auth: 2,
				},
				isErr: false,
			},
		},
		{
			name: "no such user",
			args: args{
				email: "no-email@gogin.com",
			},
			want: want{
				userAuth: nil,
				isErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userAuth, err := repo.OAuth2Login(tt.args.email)
			if (err != nil) != tt.want.isErr {
				t.Errorf("OAuth2Login() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				if err != nil {
					t.Log(err)
				}
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(userAuth, tt.want.userAuth) {
				t.Errorf("OAuth2Login() = %v, want %v", userAuth, tt.want.userAuth)
			}
		})
	}
}
