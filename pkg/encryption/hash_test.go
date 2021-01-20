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

func TestMD5_Hash(t *testing.T) {
	type args struct {
		salt1  string
		salt2  string
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
				salt1:  "hogehoge",
				salt2:  "foobar",
				target: "foo-bar-password",
			},
			want: want{
				result: "aff2d2038c7c92c6e512538c54abcd87",
			},
		},
		{
			name: "happy path 2",
			args: args{
				salt1:  "hogehogehoge",
				salt2:  "foobar",
				target: "foo-bar-password",
			},
			want: want{
				result: "4fa9c865064dadad1a06a73bf0906f7f",
			},
		},
		{
			name: "no target",
			args: args{
				salt1:  "hogehoge",
				salt2:  "foobar",
				target: "",
			},
			want: want{
				result: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md5Hash := NewMD5(tt.args.salt1, tt.args.salt2)
			got := md5Hash.Hash(tt.args.target)
			if got != tt.want.result {
				t.Errorf("md5Hash.Hash(): got = %s, want %s", got, tt.want.result)
			}
		})
	}
}

func TestMD5_HashWith(t *testing.T) {
	type args struct {
		salt1      string
		salt2      string
		additional string
		target     string
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
				salt1:      "hogehoge",
				salt2:      "foobar",
				additional: "plus-abcdefghijk",
				target:     "foo-bar-password",
			},
			want: want{
				result: "981b03f25cd905d4b7e48e225a633cae",
			},
		},
		{
			name: "happy path 2",
			args: args{
				salt1:      "hogehoge",
				salt2:      "foobar",
				additional: "plus-xxxxxx",
				target:     "foo-bar-password",
			},
			want: want{
				result: "25a90a1eda5fe72aee7b675a68538954",
			},
		},
		{
			name: "happy path 3, no additional",
			args: args{
				salt1:      "hogehoge",
				salt2:      "foobar",
				additional: "",
				target:     "foo-bar-password",
			},
			want: want{
				result: "aff2d2038c7c92c6e512538c54abcd87",
			},
		},
		{
			name: "no target",
			args: args{
				salt1:      "hogehoge",
				salt2:      "foobar",
				additional: "plus-abcdefghijk",
				target:     "",
			},
			want: want{
				result: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md5Hash := NewMD5(tt.args.salt1, tt.args.salt2)
			got := md5Hash.HashWith(tt.args.target, tt.args.additional)
			if got != tt.want.result {
				t.Errorf("md5Hash.HashWith(): got = %s, want %s", got, tt.want.result)
			}
		})
	}
}

func TestScrypt_Hash(t *testing.T) {
	type args struct {
		salt   string
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
				salt:   "hogehoge",
				target: "foo-bar-password",
			},
			want: want{
				result: "ocOpDbQBP2aNAlZtt/56ukYUaFZh/qQHMul+BkNMesY=",
			},
		},
		{
			name: "happy path 2",
			args: args{
				salt:   "hogehogehoge",
				target: "foo-bar-password",
			},
			want: want{
				result: "CZLfzQfzjwBBwCOqeObobg8CgfkJpw1r6LXPSXgo558=",
			},
		},
		{
			name: "no target",
			args: args{
				salt:   "hogehoge",
				target: "",
			},
			want: want{
				result: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scryptHash := NewScrypt(tt.args.salt)
			got := scryptHash.Hash(tt.args.target)
			if got != tt.want.result {
				t.Errorf("scryptHash.Hash(): got = %s, want %s", got, tt.want.result)
			}
		})
	}
}
