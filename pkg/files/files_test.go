package files

import (
	"testing"
)

func TestIsStaticFile(t *testing.T) {
	type args struct {
		fileName string
	}
	type want struct {
		result bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				fileName: "foo-bar.go",
			},
			want: want{
				result: true,
			},
		},
		{
			name: "happy path 2",
			args: args{
				fileName: "foo-bar.go.txt",
			},
			want: want{
				result: true,
			},
		},
		{
			name: "no extension",
			args: args{
				fileName: "foo-bar",
			},
			want: want{
				result: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsStaticFile(tt.args.fileName)
			if got != tt.want.result {
				t.Errorf("IsStaticFile(): got = %t, want %t", got, tt.want.result)
			}
		})
	}
}

func TestIsInvisiblefile(t *testing.T) {
	type args struct {
		fileName string
	}
	type want struct {
		result bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				fileName: ".foo-bar",
			},
			want: want{
				result: true,
			},
		},
		{
			name: "happy path 2",
			args: args{
				fileName: ".foo-bar.go",
			},
			want: want{
				result: true,
			},
		},
		{
			name: "no dot on the top",
			args: args{
				fileName: "foo-bar",
			},
			want: want{
				result: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsInvisiblefile(tt.args.fileName)
			if got != tt.want.result {
				t.Errorf("IsInvisiblefile(): got = %t, want %t", got, tt.want.result)
			}
		})
	}
}

func TestIsExtFile(t *testing.T) {
	type args struct {
		fileName  string
		extension string
	}
	type want struct {
		result bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				fileName:  "foo-bar.txt",
				extension: "txt",
			},
			want: want{
				result: true,
			},
		},
		{
			name: "happy path 2",
			args: args{
				fileName:  ".foo-bar.go",
				extension: "go",
			},
			want: want{
				result: true,
			},
		},
		{
			name: "no extension",
			args: args{
				fileName:  "foo-bar",
				extension: "txt",
			},
			want: want{
				result: false,
			},
		},
		{
			name: "no extension at second parameter",
			args: args{
				fileName:  "foo-bar.txt",
				extension: "",
			},
			want: want{
				result: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsExtFile(tt.args.fileName, tt.args.extension)
			if got != tt.want.result {
				t.Errorf("IsExtFile(): got = %t, want %t", got, tt.want.result)
			}
		})
	}
}
