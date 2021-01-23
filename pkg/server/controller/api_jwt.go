package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	jsonresp "github.com/hiromaily/go-gin-wrapper/pkg/server/response/json"
)

// APIJWTer interface
type APIJWTer interface {
	APIJWTIndexPostAction(c *gin.Context)
}

// APIJWTIndexPostAction is JWT endpoint [POST]
func (ctl *controller) APIJWTIndexPostAction(c *gin.Context) {
	ctl.logger.Debug("controler APIJWTIndexPostAction")

	// check login
	userID, mail, err := ctl.checkAPILogin(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, err := ctl.jwter.CreateBasicToken(
		time.Now().Add(time.Minute*60).Unix(),
		strconv.Itoa(userID),
		mail,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctl.logger.Debug("APIJWTIndexPostAction", zap.String("token", token))

	// json response
	jsonresp.ResponseUserJSON(c, ctl.logger, ctl.cors, http.StatusOK, jsonresp.CreateJWTJson(token))
}
