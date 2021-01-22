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
	conf, err := NewConfig(fmt.Sprintf("%s/%s", basePath, fileName), false)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func GetEnvConf() (*Root, error) {
	conf, err := NewConfig(os.Getenv("GOGIN_CONF"), false)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
