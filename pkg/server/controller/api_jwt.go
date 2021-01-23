package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	hh "github.com/hiromaily/go-gin-wrapper/pkg/server/httpheader"
	jsonresp "github.com/hiromaily/go-gin-wrapper/pkg/server/response/json"
)

// APIJWTer interface
type APIJWTer interface {
	APIJWTIndexPostAction(ctx *gin.Context)
}

// APIJWTIndexPostAction is JWT endpoint [POST]
func (ctl *controller) APIJWTIndexPostAction(ctx *gin.Context) {
	ctl.logger.Debug("controler APIJWTIndexPostAction")

	// check login
	userID, mail, err := ctl.checkAPILogin(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	token, err := ctl.jwter.CreateBasicToken(
		time.Now().Add(time.Minute*60).Unix(),
		strconv.Itoa(userID),
		mail,
	)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctl.logger.Debug("APIJWTIndexPostAction", zap.String("token", token))

	hh.SetResponseHeader(ctx, ctl.logger)
	// FIXME
	//if ctl.corsConf.Enabled && ctx.Request.Method == "GET" {
	//	cors.SetHeader(ctx)
	//}
	// json response
	jsonresp.ResponseUserJSON(ctx, http.StatusOK, jsonresp.CreateJWTJson(token))
}
