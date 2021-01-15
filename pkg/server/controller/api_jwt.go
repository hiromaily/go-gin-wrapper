package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/auth/jwt"
	js "github.com/hiromaily/go-gin-wrapper/pkg/json"
	jsonresp "github.com/hiromaily/go-gin-wrapper/pkg/server/response/json"
)

// APIJWTIndexPostAction is JWT End Point [POST]
func (ctl *Controller) APIJWTIndexPostAction(c *gin.Context) {
	ctl.logger.Debug("APIJWTIndexPostAction")

	// login
	// check login
	userID, mail, err := ctl.CheckLoginOnAPI(c)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	ti := time.Now().Add(time.Minute * 60).Unix()
	token, err := jwt.CreateBasicToken(ti, strconv.Itoa(userID), mail)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	ctl.logger.Debug("", zap.String("token: %s", token))

	// Make json for response and return
	jsonresp.ResponseUserJSON(c, ctl.logger, ctl.cors, 0, js.CreateJWTJson(token))
}
