package csrf

import (
	"fmt"
	lg "github.com/hiromaily/golibs/log"
	//"github.com/josephspurrier/csrfbanana"
	"crypto/md5"
	"io"
	"strconv"
	"time"
)

//TODO: set on TOML
// tokenSalt is key for md5 encryption
const tokenSalt string = "goginwebservertoken"

// CreateToken is to create token string
func CreateToken() string {
	lg.Info("[CreateToken]")
	//TODO: it may be ok to return just session key.
	h := md5.New()

	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, tokenSalt)

	token := fmt.Sprintf("%x", h.Sum(nil))

	lg.Debugf("Token: %s", token)

	return token
}
