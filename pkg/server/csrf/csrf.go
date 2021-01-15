package csrf

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// TODO: set on TOML
// tokenSalt is key for md5 encryption
const tokenSalt string = "goginwebservertoken"

// CreateToken is to create token string
func CreateToken(logger *zap.Logger) string {
	logger.Info("CreateToken")
	// TODO: it may be ok to return just session key.
	h := md5.New()

	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, tokenSalt)

	token := fmt.Sprintf("%x", h.Sum(nil))

	logger.Debug("CreateToken", zap.String("token", token))

	return token
}
