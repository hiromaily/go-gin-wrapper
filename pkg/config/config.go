package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-gin-wrapper/pkg/encryption"
)

// NewConfig returns *Root config
func NewConfig(fileName string, isEncrypted bool) (*Root, error) {
	conf, err := loadConfig(fileName)
	if err != nil {
		return nil, err
	}

	if isEncrypted {
		conf.decrypt()
	}
	return conf, err
}

// load config file
func loadConfig(fileName string) (*Root, error) {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read file: %s", fileName)
	}

	var root Root
	_, err = toml.Decode(string(d), &root)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to parse: %s", fileName)
	}

	// check validation of config
	if err = root.validate(); err != nil {
		return nil, err
	}

	return &root, nil
}

func (c *Root) validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// decrypt decrypts encrypted values in config file
func (c *Root) decrypt() {
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
