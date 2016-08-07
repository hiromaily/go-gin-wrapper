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

const TOKEN_SALT string = "goginwebservertoken"

func CreateToken() string {
	//TODO: it may be ok to return just session key.
	h := md5.New()

	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	io.WriteString(h, TOKEN_SALT)

	token := fmt.Sprintf("%x", h.Sum(nil))

	lg.Debugf("Token: %s", token)

	return token
}
