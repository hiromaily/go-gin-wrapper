package token

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	type args struct {
		salt string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy path",
			args: args{
				salt: "hogehoge",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewGenerator(tt.args.salt)
			got1 := gen.Generate()
			got2 := gen.Generate()
			if got1 == "" {
				t.Error("token.Generate(): got blank")
			}
			if got1 == got2 {
				t.Error("token.Generate(): result must be changed every time")
			}
			t.Log(got1, got2)
		})
	}
}
