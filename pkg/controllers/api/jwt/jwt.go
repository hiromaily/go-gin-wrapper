package jwt

import (
	js "github.com/hiromaily/go-gin-wrapper/pkg/jsons"
	"github.com/hiromaily/go-gin-wrapper/pkg/libs/login"
	jslib "github.com/hiromaily/go-gin-wrapper/pkg/libs/response/json"
	"github.com/hiromaily/golibs/auth/jwt"
	lg "github.com/hiromaily/golibs/log"
	u "github.com/hiromaily/golibs/utils"

	//gin "gopkg.in/gin-gonic/gin.v1"
	"time"

	"github.com/gin-gonic/gin"
)

// IndexPostAction is JWT End Point [POST]
func IndexPostAction(c *gin.Context) {
	lg.Debug("[POST] IndexAction")

	//login
	//check login
	userID, mail, err := login.CheckLoginOnAPI(c)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	ti := time.Now().Add(time.Minute * 60).Unix()
	token, err := jwt.CreateBasicToken(ti, u.Itoa(userID), mail)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	lg.Debugf("token: %s", token)

	//Make json for response and return
	jslib.RtnUserJSON(c, 0, js.CreateJWTJson(token))
	return
}
