package jwt

import (
	"github.com/gin-gonic/gin"
	js "github.com/hiromaily/go-gin-wrapper/jsons"
	jslib "github.com/hiromaily/go-gin-wrapper/libs/json"
	"github.com/hiromaily/go-gin-wrapper/libs/login"
	"github.com/hiromaily/golibs/auth/jwt"
	lg "github.com/hiromaily/golibs/log"
	u "github.com/hiromaily/golibs/utils"
	"time"
)

// JWT End Point [POST]
func IndexAction(c *gin.Context) {
	lg.Debug("[POST] IndexAction")

	//login
	//check login
	userId, mail, err := login.CheckLoginAPI(c)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	ti := time.Now().Add(time.Minute * 60).Unix()
	token, err := jwt.CreateBasicToken(ti, u.Itoa(userId), mail)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	lg.Debugf("token: %s", token)

	//Make json for response and return
	jslib.RtnUserJson(c, 0, js.CreateJWTJson(token))
	return
}
