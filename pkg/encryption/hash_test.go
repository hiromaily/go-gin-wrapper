package encryption

import "testing"

func TestHashMD5(t *testing.T) {
	type args struct {
		target string
	}
	type want struct {
		result string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				target: "foo-bar-password",
			},
			want: want{
				result: "9dc28a6aac70408f93f7a6998b0b1e3f",
			},
		},
		{
			name: "no target",
			args: args{
				target: "",
			},
			want: want{
				result: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashMD5(tt.args.target)
			if got != tt.want.result {
				t.Errorf("HashMD5(): got = %s, want %s", got, tt.want.result)
			}
		})
	}
}

func TestHashSHA1(t *testing.T) {
	type args struct {
		target string
	}
	type want struct {
		result string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				target: "foo-bar-password",
			},
			want: want{
				result: "ae19f73840b1dfe27d2c69960314c5bef3c1b827",
			},
		},
		{
			name: "no target",
			args: args{
				target: "",
			},
			want: want{
				result: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashSHA1(tt.args.target)
			if got != tt.want.result {
				t.Errorf("HashSHA1(): got = %s, want %s", got, tt.want.result)
			}
		})
	}
}

func TestHashSHA256(t *testing.T) {
	type args struct {
		target string
	}
	type want struct {
		result string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				target: "foo-bar-password",
			},
			want: want{
				result: "827afe756c693ce4c1e6539f7439e41994385adb5dd4e2d0bff05169972dee19",
			},
		},
		{
			name: "no target",
			args: args{
				target: "",
			},
			want: want{
				result: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashSHA256(tt.args.target)
			if got != tt.want.result {
				t.Errorf("HashSHA256(): got = %s, want %s", got, tt.want.result)
			}
		})
	}
}
