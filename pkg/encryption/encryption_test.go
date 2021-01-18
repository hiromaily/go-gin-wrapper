package encryption

import (
	"testing"
)

func TestEncryption(t *testing.T) {
	type args struct {
		key    string
		iv     string
		target string
	}
	type want struct {
		isErr     bool
		encrypted string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				key:    "8#75aaR+ba5Ztest",
				iv:     "@~wp-7hPs<mEtest",
				target: "abcdefg@gmail.com",
			},
			want: want{
				isErr:     false,
				encrypted: "IybYFA5wpMNHrwVVkN0kjKyayumAozHV8WQoGYbQ8oo=",
			},
		},
		{
			name: "blank key",
			args: args{
				key:    "",
				iv:     "@~wp-7hPs<mEtest",
				target: "abcdefg@gmail.com",
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "blank iv",
			args: args{
				key:    "8#75aaR+ba5Ztest",
				iv:     "",
				target: "abcdefg@gmail.com",
			},
			want: want{
				isErr: true,
			},
		},
		{
			name: "blank target",
			args: args{
				key:    "8#75aaR+ba5Ztest",
				iv:     "@~wp-7hPs<mEtest",
				target: "",
			},
			want: want{
				isErr:     false,
				encrypted: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crypt, err := NewCrypt(tt.args.key, tt.args.iv)
			if (err != nil) != tt.want.isErr {
				t.Errorf("NewCrypt() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				return
			}
			if err != nil {
				return
			}
			got := crypt.EncryptBase64(tt.args.target)
			if got != tt.want.encrypted {
				t.Errorf("crypt.EncryptBase64() = %s, want %s", got, tt.want.encrypted)
				return
			}
			got, _ = crypt.DecryptBase64(got)
			if got != tt.args.target {
				t.Errorf("crypt.DecryptBase64() = %s, want %s", got, tt.args.target)
				return
			}
		})
	}
}
