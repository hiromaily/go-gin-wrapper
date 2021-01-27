package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/server/ginctx"
	"github.com/hiromaily/go-gin-wrapper/pkg/server/validator"
)

// Loginer interface
type Loginer interface {
	login(ctx *gin.Context) (int, *LoginRequest, []string)
	apiLogin(ctx *gin.Context) (int, string, error)
}

// LoginRequest is login payload
type LoginRequest struct {
	Email string `valid:"nonempty,email,min=5,max=40" field:"email" dispName:"E-Mail"`
	Pass  string `valid:"nonempty,min=8,max=20" field:"pass" dispName:"Password"`
}

// errFormat is required by validation package
var errFormat = map[string]string{
	"nonempty": "Empty is not allowed on %s",
	"email":    "Format of %s is invalid",
	"alphanum": "Only alphabet is allowd on %s",
	"min":      "At least %s of characters is required on %s",
	"max":      "At a maximum %s of characters is allowed on %s",
}

// login validates for login
func (ctl *controller) login(ctx *gin.Context) (int, *LoginRequest, []string) {
	loginRequest := &LoginRequest{
		Email: ctx.PostForm("inputEmail"),
		Pass:  ctx.PostForm("inputPassword"),
	}

	result := validator.Validate(loginRequest, false)
	if len(result) != 0 {
		errs := validator.ConvertErrorMsgs(result, errFormat)
		return 0, loginRequest, errs
	}

	userID, err := ctl.userRepo.Login(loginRequest.Email, loginRequest.Pass)
	if err != nil {
		ctl.logger.Debug("login_error", zap.Error(err))
		errs := []string{"E-mail or Password is invalid"}
		return 0, loginRequest, errs
	}
	return userID, nil, nil
}

// apiLogin validates for API login
func (ctl *controller) apiLogin(ctx *gin.Context) (int, string, error) {
	email := ctx.PostForm("inputEmail")
	loginRequest := &LoginRequest{
		Email: email,
		Pass:  ctx.PostForm("inputPassword"),
	}

	result := validator.Validate(loginRequest, false)
	if len(result) != 0 {
		return 0, "", errors.New("validation error")
	}

	userID, err := ctl.userRepo.Login(loginRequest.Email, loginRequest.Pass)
	if err != nil {
		return 0, "", errors.New("login error")
	}
	return userID, email, nil
}

// response for login page
func (ctl *controller) loginResponse(ctx *gin.Context, input *LoginRequest, msg string, errs []string) {
	token := ctl.session.GenerateToken()
	ctl.session.SetToken(ctx, token)
	ctl.logger.Debug("loginResponse",
		zap.String("session.GenerateToken()", token),
		zap.String("SetToken()", ""),
	)
	ginctx.DebugContext(ctx, ctl.logger)

	// Google/Facebook Open ID
	gURL := "/oauth2/google/signin"
	fURL := "/oauth2/facebook/signin"

	if msg == "" {
		msg = "enter"
	}
	if input == nil {
		input = &LoginRequest{}
	}

	// view
	ctx.HTML(http.StatusOK, "pages/bases/login.tmpl", gin.H{
		"message":               msg,
		"input":                 input,
		"github.com/pkg/errors": errs,
		"gintoken":              token,
		"gURL":                  gURL,
		"fURL":                  fURL,
	})
}
