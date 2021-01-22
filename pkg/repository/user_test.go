// +build integration

package repository

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/wacul/ptr"

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

var (
	conf     *config.Root
	userRepo UserRepository
	crypt    encryption.Hasher
)

func getConf(t *testing.T) *config.Root {
	if conf == nil {
		var err error
		conf, err = config.GetConf("settings.toml")
		if err != nil {
			t.Fatal(err)
		}
	}
	return conf
}

func getCrypt(t *testing.T) encryption.Hasher {
	if crypt == nil {
		crypt = encryption.NewMD5(getConf(t).Hash.Salt1, getConf(t).Hash.Salt2)
	}
	return crypt
}

func getUserRepo(t *testing.T) UserRepository {
	if userRepo == nil {
		// db
		dbConn, err := mysql.NewMySQL(getConf(t).MySQL.Test)
		if err != nil {
			t.Fatal(err)
		}
		//
		userRepo = NewUserRepository(
			dbConn,
			logger.NewZapLogger(getConf(t).Logger),
			getCrypt(t),
		)
	}
	return userRepo
}

func init() {
	boil.DebugMode = true
}

func TestLogin(t *testing.T) {
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
			userID, err := repo.Login(tt.args.email, tt.args.password)
			if (err != nil) != tt.want.isErr {
				t.Errorf("Login() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				if err != nil {
					t.Log(err)
				}
				return
			}
			if err != nil {
				return
			}
			if userID != tt.want.userID {
				t.Errorf("Login(): userID = %d, want %d", userID, tt.want.userID)
			}
		})
	}
}

func TestOAuth2Login(t *testing.T) {
	repo := getUserRepo(t)

	type args struct {
		email string
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
				isErr:    true,
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

func TestGetUserIDs(t *testing.T) {
	repo := getUserRepo(t)

	type want struct {
		len   int
		isErr bool
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "happy path 1",
			want: want{
				len:   4,
				isErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids, err := repo.GetUserIDs()
			if (err != nil) != tt.want.isErr {
				t.Errorf("GetUserIDs() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				if err != nil {
					t.Log(err)
				}
				return
			}
			if err != nil {
				return
			}
			if len(ids) != tt.want.len {
				t.Errorf("GetUserIDs() ids length = %d, want %d", len(ids), tt.want.len)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	repo := getUserRepo(t)

	type args struct {
		id string
	}
	type want struct {
		users []*user.User
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				id: "1",
			},
			want: want{
				users: []*user.User{
					{
						ID:        1,
						FirstName: "foo",
						LastName:  "bar",
						Email:     "foobar@gogin.com",
						Password:  "baa62a499e9b21940c2d763f58a25647",
						OAuth2:    0,
						Created:   ptr.Time(time.Date(2021, time.January, 10, 21, 43, 15, 0, time.UTC)),
						Updated:   ptr.Time(time.Date(2021, time.January, 10, 21, 43, 15, 0, time.UTC)),
					},
				},
				isErr: false,
			},
		},
		{
			name: "happy path : multiple user response",
			args: args{
				id: "",
			},
			want: want{
				users: []*user.User{
					{
						ID:        1,
						FirstName: "foo",
						LastName:  "bar",
						Email:     "foobar@gogin.com",
						Password:  "baa62a499e9b21940c2d763f58a25647",
						OAuth2:    0,
						Created:   ptr.Time(time.Date(2021, time.January, 10, 21, 43, 15, 0, time.UTC)),
						Updated:   ptr.Time(time.Date(2021, time.January, 10, 21, 43, 15, 0, time.UTC)),
					},
					{
						ID:        2,
						FirstName: "mark",
						LastName:  "harry",
						Email:     "mark@gogin.com",
						Password:  "d978eb967fbe04345371478a97f3c903",
						OAuth2:    0,
						Created:   ptr.Time(time.Date(2021, time.January, 11, 20, 20, 28, 0, time.UTC)),
						Updated:   ptr.Time(time.Date(2021, time.January, 11, 20, 20, 28, 0, time.UTC)),
					},
					{
						ID:        3,
						FirstName: "check",
						LastName:  "authtype1",
						Email:     "chk-authtype1@gogin.com",
						Password:  "d978eb967fbe04345371478a97f3c903",
						OAuth2:    1,
						Created:   ptr.Time(time.Date(2021, time.January, 11, 20, 20, 28, 0, time.UTC)),
						Updated:   ptr.Time(time.Date(2021, time.January, 11, 20, 20, 28, 0, time.UTC)),
					},
					{
						ID:        4,
						FirstName: "check",
						LastName:  "authtype2",
						Email:     "chk-authtype2@gogin.com",
						Password:  "d978eb967fbe04345371478a97f3c903",
						OAuth2:    2,
						Created:   ptr.Time(time.Date(2021, time.January, 11, 20, 20, 28, 0, time.UTC)),
						Updated:   ptr.Time(time.Date(2021, time.January, 11, 20, 20, 28, 0, time.UTC)),
					},
				},
				isErr: false,
			},
		},
		{
			name: "no such id user",
			args: args{
				id: "999999",
			},
			want: want{
				users: nil,
				isErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := repo.GetUsers(tt.args.id)
			if (err != nil) != tt.want.isErr {
				t.Errorf("GetUsers() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				if err != nil {
					t.Log(err)
				}
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(users, tt.want.users) {
				t.Errorf("GetUsers() = %v, want %v", users, tt.want.users)
			}
			// for only debug
			//for idx, u := range users {
			//	if !reflect.DeepEqual(u, tt.want.users[idx]) {
			//		t.Errorf("GetUsers() [%d] = %v, want %v", idx, users, tt.want.users)
			//	}
			//	t.Log(u)
			//	t.Log(tt.want.users[idx])
			//}
		})
	}
}

// db data should be reset by `make setup-testdb`
func TestInsertUser(t *testing.T) {
	repo := getUserRepo(t)

	user1 := &user.User{
		FirstName: "first",
		LastName:  "last",
		Email:     "test@gogin-test.com",
		Password:  "plain-text",
		OAuth2:    1,
	}

	oldTimeUnix := time.Now().UTC().UnixNano()

	type args struct {
		user *user.User
	}
	type want struct {
		id    int
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				user: user1,
			},
			want: want{
				id:    6,
				isErr: false,
			},
		},
		{
			name: "insert duplicated user",
			args: args{
				user: user1,
			},
			want: want{
				isErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.InsertUser(tt.args.user)
			if (err != nil) != tt.want.isErr {
				t.Errorf("InsertUser() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				if err != nil {
					t.Log(err)
				}
				return
			}
			if err != nil {
				return
			}
			if id != tt.want.id {
				t.Errorf("InsertUser() id = %d, want %d", id, tt.want.id)
			}
			// select
			users, err := repo.GetUsers(strconv.Itoa(id))
			if err != nil {
				t.Fatal(err)
			}
			// validate
			if tt.args.user.FirstName != users[0].FirstName {
				t.Errorf("InsertUser() FirstName = %s, want %s", users[0].FirstName, tt.args.user.FirstName)
			}
			if tt.args.user.LastName != users[0].LastName {
				t.Errorf("InsertUser() LastName = %s, want %s", users[0].LastName, tt.args.user.LastName)
			}
			if tt.args.user.Email != users[0].Email {
				t.Errorf("InsertUser() Email = %s, want %s", users[0].Email, tt.args.user.Email)
			}
			if getCrypt(t).Hash(tt.args.user.Password) != users[0].Password {
				t.Errorf("InsertUser() Password = %s, want %s", users[0].Password, getCrypt(t).Hash(tt.args.user.Password))
			}
			if tt.args.user.OAuth2 != users[0].OAuth2 {
				t.Errorf("InsertUser() OAuth2 = %d, want %d", users[0].OAuth2, tt.args.user.OAuth2)
			}
			if users[0].Created.UnixNano() < oldTimeUnix {
				t.Errorf("InsertUser() Created = %d, old time unit(): %d", users[0].Created.UnixNano(), oldTimeUnix)
			}
			if users[0].Updated.UnixNano() < oldTimeUnix {
				t.Errorf("InsertUser() Updated = %d, old time unit(): %d", users[0].Updated.UnixNano(), oldTimeUnix)
			}
			// debug
			t.Log(users[0])
		})
	}
}

// db data should be reset by `make setup-testdb`
func TestUpdateUser(t *testing.T) {
	repo := getUserRepo(t)

	user1 := &user.User{
		FirstName: "first-updated",
		LastName:  "last-updated",
		Email:     "test-updated@gogin-test.com",
		Password:  "plain-text-updated",
	}

	type args struct {
		user   *user.User
		userID int
	}
	type want struct {
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				user:   user1,
				userID: 1,
			},
			want: want{
				isErr: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.UpdateUser(tt.args.user, tt.args.userID)
			if (err != nil) != tt.want.isErr {
				t.Errorf("UpdateUser() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				if err != nil {
					t.Log(err)
				}
				return
			}
			if err != nil {
				return
			}
			// select
			users, err := repo.GetUsers(strconv.Itoa(tt.args.userID))
			if err != nil {
				t.Fatal(err)
			}
			// validate
			if tt.args.user.FirstName != users[0].FirstName {
				t.Errorf("UpdateUser() FirstName = %s, want %s", users[0].FirstName, tt.args.user.FirstName)
			}
			if tt.args.user.LastName != users[0].LastName {
				t.Errorf("UpdateUser() LastName = %s, want %s", users[0].LastName, tt.args.user.LastName)
			}
			if tt.args.user.Email != users[0].Email {
				t.Errorf("UpdateUser() Email = %s, want %s", users[0].Email, tt.args.user.Email)
			}
			if getCrypt(t).Hash(tt.args.user.Password) != users[0].Password {
				t.Errorf("UpdateUser() Password = %s, want %s", users[0].Password, getCrypt(t).Hash(tt.args.user.Password))
			}
			//if users[0].Updated.UnixNano() < oldTimeUnix {
			//	t.Errorf("InsertUser() Updated = %d, old time unit(): %d", users[0].Updated.UnixNano(), oldTimeUnix)
			//}
			// debug
			t.Log(users[0])
		})
	}
}
