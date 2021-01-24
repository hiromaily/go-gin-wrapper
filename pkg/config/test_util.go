package config

import (
	"fmt"
	"os"
)

// GetConf returns *Root for unittest
// e.g. `conf, err = config.GetConf("settings.toml")`
func GetConf(fileName string) (*Root, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	basePath := fmt.Sprintf("%s/../../configs", pwd)

	return newConf(fmt.Sprintf("%s/%s", basePath, fileName))
}

// GetEnvConf returns *Root from environment variable `$GO-GIN_CONF` for unittest
func GetEnvConf() (*Root, error) {
	return newConf(os.Getenv("GO_GIN_CONF"))
}

func newConf(filePath string) (*Root, error) {
	conf, err := NewConfig(filePath, false)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
