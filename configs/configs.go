package configs

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	enc "github.com/hiromaily/golibs/cipher/encryption"
	u "github.com/hiromaily/golibs/utils"
	"io/ioutil"
	"os"
)

/* singleton */
var (
	conf         *Config
	tomlFileName = "./configs/settings.toml"
)

// Config is of root
type Config struct {
	Environment string
	Server      *ServerConfig
	Proxy       ProxyConfig
	Auth        *AuthConfig
	MySQL       *MySQLConfig
	Redis       *RedisConfig
	Mongo       *MongoConfig `toml:"mongodb"`
	Aws         *AwsConfig
	Develop     DevelopConfig
}

// ServerConfig is for web server
type ServerConfig struct {
	Scheme    string          `toml:"scheme"`
	Host      string          `toml:"host"`
	Port      int             `toml:"port"`
	Docs      DocsConfig      `toml:"docs"`
	Log       LogConfig       `toml:"log"`
	Session   SessionConfig   `toml:"session"`
	BasicAuth BasicAuthConfig `toml:"basic_auth"`
}

// DocsConfig is path for document root of webserver
type DocsConfig struct {
	Path string `toml:"path"`
}

// LogConfig is for Log
type LogConfig struct {
	Level uint8  `toml:"level"`
	Path  string `toml:"path"`
}

// SessionConfig is for Session
type SessionConfig struct {
	Name     string `toml:"name"`
	Key      string `toml:"key"`
	MaxAge   int    `toml:"max_age"`
	Secure   bool   `toml:"secure"`
	HTTPOnly bool   `toml:"http_only"`
}

// BasicAuthConfig is for Basic Auth
type BasicAuthConfig struct {
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

// ProxyConfig is for base of Reverse Proxy Server
type ProxyConfig struct {
	Mode   uint8             `toml:"mode"` //0:off, 1:go, 2,nginx
	Server ProxyServerConfig `toml:"server"`
}

// ProxyServerConfig is for Reverse Proxy Server
type ProxyServerConfig struct {
	Scheme  string    `toml:"scheme"`
	Host    string    `toml:"host"`
	Port    int       `toml:"port"`
	WebPort []int     `toml:"web_port"`
	Log     LogConfig `toml:"log"`
}

// AuthConfig is for Auth
type AuthConfig struct {
	API      *APIConfig      `toml:"api"`
	JWT      *JWTConfig      `toml:"jwt"`
	Google   *GoogleConfig   `toml:"google"`
	Facebook *FacebookConfig `toml:"facebook"`
}

// APIConfig is for Rest API
type APIConfig struct {
	Header string `toml:"header"`
	Key    string `toml:"key"`
	Ajax   bool   `toml:"only_ajax"`
}

// JWTConfig is for JWT Auth
type JWTConfig struct {
	Mode       uint8  `toml:"mode"`
	Secret     string `toml:"secret_code"`
	PrivateKey string `toml:"private_key"`
	PublicKey  string `toml:"public_key"`
}

// GoogleConfig is for Google OAuth2
type GoogleConfig struct {
	Encrypted    bool   `toml:"encrypted"`
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	CallbackURL  string `toml:"call_back_url"`
}

// FacebookConfig is for Facebook OAuth2
type FacebookConfig struct {
	Encrypted    bool   `toml:"encrypted"`
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	CallbackURL  string `toml:"call_back_url"`
}

// MySQLConfig is for MySQL Server
type MySQLConfig struct {
	*MySQLContentConfig
	Test *MySQLContentConfig `toml:"test"`
}

// MySQLContentConfig is for MySQL Server
type MySQLContentConfig struct {
	Encrypted bool   `toml:"encrypted"`
	Host      string `toml:"host"`
	Port      uint16 `toml:"port"`
	DbName    string `toml:"dbname"`
	User      string `toml:"user"`
	Pass      string `toml:"pass"`
}

// RedisConfig is for Redis Server
type RedisConfig struct {
	Encrypted bool   `toml:"encrypted"`
	Host      string `toml:"host"`
	Port      uint16 `toml:"port"`
	Pass      string `toml:"pass"`
	Session   bool   `toml:"session"`
}

// MongoConfig is for MongoDB Server
type MongoConfig struct {
	Encrypted bool   `toml:"encrypted"`
	Host      string `toml:"host"`
	Port      uint16 `toml:"port"`
	DbName    string `toml:"dbname"`
	User      string `toml:"user"`
	Pass      string `toml:"pass"`
}

// AwsConfig for Amazon Web Service
type AwsConfig struct {
	Encrypted bool   `toml:"encrypted"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
	Region    string `toml:"region"`
}

// DevelopConfig is for development environment
type DevelopConfig struct {
	ProfileEnable bool `toml:"profile_enable"`
	RecoverEnable bool `toml:"recover_enable"`
}

var checkTOMLKeys = [][]string{
	{"environment"},
	{"server", "scheme"},
	{"server", "host"},
	{"server", "port"},
	{"server", "docs", "path"},
	{"server", "log", "level"},
	{"server", "log", "path"},
	{"server", "session", "name"},
	{"server", "session", "key"},
	{"server", "session", "max_age"},
	{"server", "session", "secure"},
	{"server", "session", "http_only"},
	{"server", "basic_auth", "user"},
	{"server", "basic_auth", "pass"},
	{"proxy", "mode"},
	{"proxy", "server", "scheme"},
	{"proxy", "server", "host"},
	{"proxy", "server", "port"},
	{"proxy", "server", "log", "level"},
	{"proxy", "server", "log", "path"},
	{"auth", "api", "header"},
	{"auth", "api", "key"},
	{"auth", "api", "only_ajax"},
	{"auth", "jwt", "mode"},
	{"auth", "jwt", "secret_code"},
	{"auth", "jwt", "private_key"},
	{"auth", "jwt", "public_key"},
	{"auth", "google", "encrypted"},
	{"auth", "google", "client_id"},
	{"auth", "google", "client_secret"},
	{"auth", "google", "call_back_url"},
	{"auth", "facebook", "encrypted"},
	{"auth", "facebook", "client_id"},
	{"auth", "facebook", "client_secret"},
	{"auth", "facebook", "call_back_url"},
	{"mysql", "encrypted"},
	{"mysql", "host"},
	{"mysql", "port"},
	{"mysql", "dbname"},
	{"mysql", "user"},
	{"mysql", "pass"},
	{"mysql", "test", "encrypted"},
	{"mysql", "test", "host"},
	{"mysql", "test", "port"},
	{"mysql", "test", "dbname"},
	{"mysql", "test", "user"},
	{"mysql", "test", "pass"},
	{"redis", "encrypted"},
	{"redis", "host"},
	{"redis", "port"},
	{"redis", "pass"},
	{"redis", "session"},
	{"mongodb", "encrypted"},
	{"mongodb", "host"},
	{"mongodb", "port"},
	{"mongodb", "dbname"},
	{"mongodb", "user"},
	{"mongodb", "pass"},
	{"aws", "encrypted"},
	{"aws", "access_key"},
	{"aws", "secret_key"},
	{"aws", "region"},
	{"develop", "profile_enable"},
	{"develop", "recover_enable"},
}

func init() {
	tomlFileName = os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-gin-wrapper/configs/settings.toml"
}

//check validation of config
func validateConfig(conf *Config, md *toml.MetaData) error {
	//for protection when debugging on non production environment
	var errStrings []string

	//Check added new items on toml
	// environment
	//if !md.IsDefined("environment") {
	//	errStrings = append(errStrings, "environment")
	//}

	format := "[%s]"
	inValid := false
	for _, keys := range checkTOMLKeys {
		if !md.IsDefined(keys...) {
			switch len(keys) {
			case 1:
				format = "[%s]"
			case 2:
				format = "[%s] %s"
			case 3:
				format = "[%s.%s] %s"
			default:
				//invalid check string
				inValid = true
				break
			}
			keysIfc := u.SliceStrToInterface(keys)
			errStrings = append(errStrings, fmt.Sprintf(format, keysIfc...))
		}
	}

	// Error
	if inValid {
		return errors.New("Error: Check Text has wrong number of parameter")
	}
	if len(errStrings) != 0 {
		return fmt.Errorf("Error: There are lacks of keys : %#v \n", errStrings)
	}

	return nil
}

// load configfile
func loadConfig(fileName string) (*Config, error) {
	if fileName != "" {
		tomlFileName = fileName
	}

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

// New is create instance
func New(fileName string) {
	var err error
	if conf == nil {
		conf, err = loadConfig(fileName)
	}
	if err != nil {
		panic(err)
	}
}

// GetConf is to get config instance
func GetConf() *Config {
	var err error
	if conf == nil {
		conf, err = loadConfig("")
	}
	if err != nil {
		panic(err)
	}

	return conf
}

// SetTOMLPath is to set toml file path
func SetTOMLPath(path string) {
	tomlFileName = path
}

// Cipher is to decrypt crypted string on config
func Cipher() {
	crypt := enc.GetCrypt()

	if conf.Auth.Google.Encrypted {
		c := conf.Auth.Google
		c.ClientID, _ = crypt.DecryptBase64(c.ClientID)
		c.ClientSecret, _ = crypt.DecryptBase64(c.ClientSecret)
	}

	if conf.Auth.Facebook.Encrypted {
		c := conf.Auth.Facebook
		c.ClientID, _ = crypt.DecryptBase64(c.ClientID)
		c.ClientSecret, _ = crypt.DecryptBase64(c.ClientSecret)
	}

	if conf.MySQL.Encrypted {
		c := conf.MySQL
		c.Host, _ = crypt.DecryptBase64(c.Host)
		c.DbName, _ = crypt.DecryptBase64(c.DbName)
		c.User, _ = crypt.DecryptBase64(c.User)
		c.Pass, _ = crypt.DecryptBase64(c.Pass)
	}

	if conf.MySQL.Test.Encrypted {
		c := conf.MySQL.Test
		c.Host, _ = crypt.DecryptBase64(c.Host)
		c.DbName, _ = crypt.DecryptBase64(c.DbName)
		c.User, _ = crypt.DecryptBase64(c.User)
		c.Pass, _ = crypt.DecryptBase64(c.Pass)
	}

	if conf.Redis.Encrypted {
		c := conf.Redis
		c.Host, _ = crypt.DecryptBase64(c.Host)
		c.Pass, _ = crypt.DecryptBase64(c.Pass)
	}

	if conf.Mongo.Encrypted {
		c := conf.Mongo
		c.Host, _ = crypt.DecryptBase64(c.Host)
		c.DbName, _ = crypt.DecryptBase64(c.DbName)
		c.User, _ = crypt.DecryptBase64(c.User)
		c.Pass, _ = crypt.DecryptBase64(c.Pass)
	}

	if conf.Aws.Encrypted {
		c := conf.Aws
		c.AccessKey, _ = crypt.DecryptBase64(c.AccessKey)
		c.SecretKey, _ = crypt.DecryptBase64(c.SecretKey)
	}
}
