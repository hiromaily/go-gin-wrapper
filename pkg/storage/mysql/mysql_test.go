// +build integration

package mysql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-gin-wrapper/pkg/config"
)

// Note
// - run the below commands before this unittest
// $ docker-compose mysql
// $ make setup-testdb

func TestNewMySQL(t *testing.T) {
	type args struct {
		conf *config.MySQLContent
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
			// see ./configs/settings.toml
			// [mysql.test]
			name: "happy path 1",
			args: args{
				conf: &config.MySQLContent{
					Host:   "127.0.0.1",
					Port:   13306,
					DBName: "go-gin-test",
					User:   "guestuser",
					Pass:   "secret123",
				},
			},
			want: want{
				isErr: false,
			},
		},
		{
			name: "wrong host, timeout should happen",
			args: args{
				conf: &config.MySQLContent{
					Host:   "127.0.0.100",
					Port:   13306,
					DBName: "go-gin-test",
					User:   "guestuser",
					Pass:   "secret123",
				},
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "wrong port",
			args: args{
				conf: &config.MySQLContent{
					Host:   "127.0.0.1",
					Port:   13316,
					DBName: "go-gin-test",
					User:   "guestuser",
					Pass:   "secret123",
				},
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "wrong db name",
			args: args{
				conf: &config.MySQLContent{
					Host:   "127.0.0.1",
					Port:   13306,
					DBName: "go-gin-wrong-test",
					User:   "guestuser",
					Pass:   "secret123",
				},
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "wrong user name",
			args: args{
				conf: &config.MySQLContent{
					Host:   "127.0.0.1",
					Port:   13306,
					DBName: "go-gin-test",
					User:   "guest-wrong-user",
					Pass:   "secret123",
				},
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "MySQLContent is nil",
			args: args{
				conf: nil,
			},
			want: want{
				isErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbConn, err := NewMySQL(tt.args.conf)
			if (err != nil) != tt.want.isErr {
				t.Errorf("NewMySQL() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				if err != nil {
					t.Log(err)
				}
				return
			}
			if err != nil {
				return
			}
			if dbConn == nil {
				t.Errorf("NewMySQL(): dbConn should NOT be nil")
			}
		})
	}
}
