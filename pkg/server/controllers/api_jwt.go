package controllers

import (
	"time"

	"github.com/gin-gonic/gin"

	js "github.com/hiromaily/go-gin-wrapper/pkg/json"
	jslib "github.com/hiromaily/go-gin-wrapper/pkg/server/response/json"
	"github.com/hiromaily/golibs/auth/jwt"
	lg "github.com/hiromaily/golibs/log"
	u "github.com/hiromaily/golibs/utils"
)

// APIJWTIndexPostAction is JWT End Point [POST]
func (ctl *Controller) APIJWTIndexPostAction(c *gin.Context) {
	lg.Debug("[POST] IndexAction")

	// login
	// check login
	userID, mail, err := ctl.CheckLoginOnAPI(c)
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

	// Make json for response and return
	jslib.ResponseUserJSON(c, ctl.cors, 0, js.CreateJWTJson(token))
}
