package files

import (
	"fmt"
	"os"
	"testing"
)

func TestFilelist(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	basePath := fmt.Sprintf("%s/../..", pwd)

	type args struct {
		filePath string
		exts     []string
	}
	type want struct {
		isErr bool
		len   int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				filePath: fmt.Sprintf("%s/web/templates/pages", basePath),
				exts:     []string{"tmpl"},
			},
			want: want{
				isErr: false,
				len:   12,
			},
		},
		{
			name: "no files",
			args: args{
				filePath: fmt.Sprintf("%s/web/templates/pages", basePath),
				exts:     []string{"hoge"},
			},
			want: want{
				isErr: false,
				len:   0,
			},
		},
		{
			name: "no directory",
			args: args{
				filePath: "/aaa/bbb",
				exts:     []string{"txt"},
			},
			want: want{
				isErr: true,
				len:   0,
			},
		},
		{
			name: "nil extensions",
			args: args{
				filePath: fmt.Sprintf("%s/web/templates/pages", basePath),
				exts:     nil,
			},
			want: want{
				isErr: true,
				len:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileList(tt.args.filePath, tt.args.exts)
			if (err != nil) != tt.want.isErr {
				t.Errorf("GetFileList() actual error: %t, want error: %t", err != nil, tt.want.isErr)
				return
			}
			if err != nil {
				return
			}
			if len(got) != tt.want.len {
				t.Errorf("GetFileList(): got length = %d, want %d", len(got), tt.want.len)
			}
		})
	}
}

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
