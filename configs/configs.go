package configs

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

/* singleton */
var conf *Config

var tomlFileName string = "./configs/settings.toml"

type Config struct {
	Environment string
	Server      ServerConfig
	MySQL       MySQLConfig
	Aws         AwsConfig
}

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type MySQLConfig struct {
	Host   string `toml:"host"`
	Port   uint16 `toml:"port"`
	DbName string `toml:"dbname"`
	User   string `toml:"user"`
	Pass   string `toml:"pass"`
}

type AwsConfig struct {
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
	Region    string `toml:"region"`
}

//check validation of config
func validateConfig(conf *Config, md *toml.MetaData) error {
	//for protection when debugging on non production environment
	var errStrings []string

	//新規追加したtomlのkeyをチェック(通常は通らないが、初期項目追加時用)
	if !md.IsDefined("environment") {
		errStrings = append(errStrings, "environment")
	}

	if !md.IsDefined("server", "host") {
		errStrings = append(errStrings, "[server] host")
	}
	if !md.IsDefined("server", "port") {
		errStrings = append(errStrings, "[server] port")
	}
	if !md.IsDefined("mysql", "host") {
		errStrings = append(errStrings, "[mysql] host")
	}
	if !md.IsDefined("mysql", "dbname") {
		errStrings = append(errStrings, "[mysql] dbname")
	}
	if !md.IsDefined("mysql", "user") {
		errStrings = append(errStrings, "[mysql] user")
	}
	if !md.IsDefined("mysql", "pass") {
		errStrings = append(errStrings, "[mysql] pass")
	}

	if len(errStrings) != 0 {
		return fmt.Errorf("Error  There are lacks of keys : %#v \n", errStrings)
	}

	return nil
}

// load configfile
func loadConfig() (*Config, error) {
	d, err := ioutil.ReadFile(tomlFileName)
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading %s: %s", tomlFileName, err)
	}

	var config Config
	md, err := toml.Decode(string(d), &config)
	if err != nil {
		return nil, fmt.Errorf(
			"Error parsing %s: %s(%v)", tomlFileName, err, md)
	}

	//check validation of config
	err = validateConfig(&config, &md)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// singleton architecture
func New() {
	var err error
	if conf == nil {
		conf, err = loadConfig()
	}
	if err != nil {
		panic(err)
	}
}

// singleton architecture
func GetConfInstance() *Config {
	var err error
	if conf == nil {
		conf, err = loadConfig()
	}
	if err != nil {
		panic(err)
	}

	return conf
}

func SetTomlPath(path string) {
	tomlFileName = path
}
