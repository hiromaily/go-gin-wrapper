package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
)

// Config is of root
type Config struct {
	IsSignal bool    `toml:"is_signal"`
	Logger   *Logger `toml:"logger"`
	Server   *ServerConfig
	Proxy    *ProxyConfig
	API      *APIConfig
	Auth     *AuthConfig
	MySQL    *MySQLConfig `toml:"mysql" validate:"required"`
	Redis    *RedisConfig `toml:"redis" validate:"required"`
	Develop  *DevelopConfig
}

// Logger logger info
type Logger struct {
	Service      string `toml:"service" validate:"required"`
	Env          string `toml:"env" validate:"oneof=dev prod custom"`
	Level        string `toml:"level" validate:"required"`
	IsStackTrace bool   `toml:"is_stacktrace"`
}

// ServerConfig is for web server
type ServerConfig struct {
	IsRelease bool            `toml:"is_release"`
	Scheme    string          `toml:"scheme" validate:"required"`
	Host      string          `toml:"host" validate:"required"`
	Port      int             `toml:"port" validate:"required"`
	Docs      DocsConfig      `toml:"docs"`
	Log       LogConfig       `toml:"log" validate:"required"`
	Session   SessionConfig   `toml:"session" validate:"required"`
	BasicAuth BasicAuthConfig `toml:"basic_auth" validate:"required"`
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
	Mode   uint8             `toml:"mode"` // 0:off, 1:go, 2,nginx
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

// APIConfig is for Rest API
type APIConfig struct {
	Ajax   bool          `toml:"only_ajax"`
	CORS   *CORSConfig   `toml:"cors"`
	Header *HeaderConfig `toml:"header"`
	JWT    *JWTConfig    `toml:"jwt"`
}

// CORSConfig is for CORS
type CORSConfig struct {
	Enabled     bool     `toml:"enabled"`
	Origins     []string `toml:"origins"`
	Headers     []string `toml:"headers"`
	Methods     []string `toml:"methods"`
	Credentials bool     `toml:"credentials"`
}

// HeaderConfig is added original header for authentication
type HeaderConfig struct {
	Enabled bool   `toml:"enabled"`
	Header  string `toml:"header"`
	Key     string `toml:"key"`
}

// JWTConfig is for JWT Auth
type JWTConfig struct {
	Mode       uint8  `toml:"mode"` // 0:off, 1:HMAC, 2:RSA
	Secret     string `toml:"secret_code"`
	PrivateKey string `toml:"private_key"`
	PublicKey  string `toml:"public_key"`
}

// AuthConfig is for Auth
type AuthConfig struct {
	Google   *GoogleConfig   `toml:"google"`
	Facebook *FacebookConfig `toml:"facebook"`
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
	Encrypted  bool   `toml:"encrypted"`
	Host       string `toml:"host"`
	Port       uint16 `toml:"port"`
	DBName     string `toml:"dbname"`
	User       string `toml:"user"`
	Pass       string `toml:"pass"`
	IsDebugLog bool   `toml:"is_debug_log"`
}

// RedisConfig is for Redis Server
type RedisConfig struct {
	Encrypted bool   `toml:"encrypted"`
	Host      string `toml:"host"`
	Port      uint16 `toml:"port"`
	Pass      string `toml:"pass"`
	IsSession bool   `toml:"is_session"`
	IsHeroku  bool   `toml:"is_heroku"`
}

// DevelopConfig is for development environment
type DevelopConfig struct {
	ProfileEnable bool `toml:"profile_enable"`
	RecoverEnable bool `toml:"recover_enable"`
}

// New is create instance
func New(fileName string, cipherFlg bool) (*Config, error) {
	conf, err := loadConfig(fileName)
	if err != nil {
		return nil, err
	}

	if cipherFlg {
		conf.Cipher()
	}
	return conf, err
}

// load configfile
func loadConfig(fileName string) (*Config, error) {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read file: %s", fileName)
	}

	var config Config
	_, err = toml.Decode(string(d), &config)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to parse: %s", fileName)
	}

	// check validation of config
	if err = config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// Cipher is to decrypt crypted string on config
func (c *Config) Cipher() {
	crypt := encryption.GetCrypt()

	if c.Auth.Google.Encrypted {
		ag := c.Auth.Google
		ag.ClientID, _ = crypt.DecryptBase64(ag.ClientID)
		ag.ClientSecret, _ = crypt.DecryptBase64(ag.ClientSecret)
	}

	if c.Auth.Facebook.Encrypted {
		ag := c.Auth.Facebook
		ag.ClientID, _ = crypt.DecryptBase64(ag.ClientID)
		ag.ClientSecret, _ = crypt.DecryptBase64(ag.ClientSecret)
	}

	if c.MySQL.Encrypted {
		m := c.MySQL
		m.Host, _ = crypt.DecryptBase64(m.Host)
		m.DBName, _ = crypt.DecryptBase64(m.DBName)
		m.User, _ = crypt.DecryptBase64(m.User)
		m.Pass, _ = crypt.DecryptBase64(m.Pass)
	}

	if c.MySQL.Test.Encrypted {
		mt := c.MySQL.Test
		mt.Host, _ = crypt.DecryptBase64(mt.Host)
		mt.DBName, _ = crypt.DecryptBase64(mt.DBName)
		mt.User, _ = crypt.DecryptBase64(mt.User)
		mt.Pass, _ = crypt.DecryptBase64(mt.Pass)
	}

	if c.Redis.Encrypted {
		r := c.Redis
		r.Host, _ = crypt.DecryptBase64(r.Host)
		r.Pass, _ = crypt.DecryptBase64(r.Pass)
	}
}
