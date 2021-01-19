package heroku

import (
	"os"
	"testing"
)

func TestGetRedisInfo(t *testing.T) {
	type args struct {
		target string
	}
	type want struct {
		isErr bool
		host  string
		pass  string
		port  int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				target: "redis://h:password@host:12345",
			},
			want: want{
				isErr: false,
				host:  "host",
				pass:  "password",
				port:  12345,
			},
		},
		{
			name: "wrong string",
			args: args{
				target: "foo-bar-foo-bar",
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "wrong string, no `//`",
			args: args{
				target: "redis:||h:password@host:12345",
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "wrong string, no first `:`",
			args: args{
				target: "redis://h-password@host:12345",
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "wrong string, no last `:`",
			args: args{
				target: "redis://h:password@host=12345",
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "wrong string, no `@`",
			args: args{
				target: "redis://h:password=host:12345",
			},
			want: want{
				isErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range []string{"", "env"} {
				target := tt.args.target
				if v == "env" {
					target = ""
					os.Setenv("REDIS_URL", tt.args.target)
				}
				host, pass, port, err := GetRedisInfo(target)
				if (err != nil) != tt.want.isErr {
					t.Errorf("GetRedisInfo() actual error: %t, want error: %t", err != nil, tt.want.isErr)
					return
				}
				if err != nil {
					return
				}
				if host != tt.want.host {
					t.Errorf("crypt.GetRedisInfo(): host = %s, want %s", host, tt.want.host)
					return
				}
				if pass != tt.want.pass {
					t.Errorf("envVariablecrypt.GetRedisInfo(): pass = %s, want %s", pass, tt.want.pass)
					return
				}
				if port != tt.want.port {
					t.Errorf("crypt.GetRedisInfo(): port = %d, want %d", port, tt.want.port)
					return
				}
			}
		})
	}
}
