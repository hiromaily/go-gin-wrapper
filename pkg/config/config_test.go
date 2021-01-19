package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/hiromaily/go-gin-wrapper/pkg/files"
)

func TestConfig(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	basePath := fmt.Sprintf("%s/../../configs", pwd)
	files, err := files.GetFileList(basePath, []string{"toml"})
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Log(file)
		if _, err := NewConfig(file, false); err != nil {
			t.Error(err)
		}
	}
}
