package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"

	"github.com/hiromaily/go-gin-wrapper/pkg/model/user"
)

// OAuther interface
type OAuther interface {
	OAuth2SignInGoogleAction(ctx *gin.Context)
	OAuth2SignInFacebookAction(ctx *gin.Context)
	OAuth2CallbackGoogleAction(ctx *gin.Context)
	OAuth2CallbackFacebookAction(ctx *gin.Context)
}

// ResGoogle is for response data from google
type ResGoogle struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`        // full name
	FirstName     string `json:"given_name"`  // first name
	LastName      string `json:"family_name"` // last name
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

// ResFacebook is for response data from facebook
type ResFacebook struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified"`
	Name          string `json:"name"`       // full name
	FirstName     string `json:"first_name"` // first name
	LastName      string `json:"last_name"`  // last name
	Link          string `json:"link"`
	Picture       FBPic  `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

// FBPic is structure of picture on facebook
type FBPic struct {
	Data struct {
		IsSilhouette bool   `json:"is_silhouette"`
		URL          string `json:"url"`
	}
}

const (
	// GoogleAuth is for google
	GoogleAuth uint8 = 1
	// FacebookAuth is for Facebook
	FacebookAuth uint8 = 2
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "",
		ClientID:     "",
		ClientSecret: "",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	facebookOauthConfig = &oauth2.Config{
		RedirectURL:  "",
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"public_profile", "email"},
		// Scopes:   []string{"first_name", "last_name", "link", "picture", "email"},
		Endpoint: facebook.Endpoint,
	}
)

// OAuth2SignInGoogleAction is sign in by Google [GET]
func (ctl *controller) OAuth2SignInGoogleAction(ctx *gin.Context) {
	ctl.logger.Info("OAuth2SignInGoogleAction")

	auth := ctl.authConf.Google

	googleOauthConfig.RedirectURL = auth.CallbackURL
	googleOauthConfig.ClientID = auth.ClientID
	googleOauthConfig.ClientSecret = auth.ClientSecret

	// token
	token := ctl.session.GenerateToken()
	ctl.session.SetToken(ctx, token)

	url := googleOauthConfig.AuthCodeURL(token)
	ctx.Redirect(http.StatusTemporaryRedirect, url) // 307
}

// OAuth2SignInFacebookAction is sign in by Facebook [GET]
func (ctl *controller) OAuth2SignInFacebookAction(ctx *gin.Context) {
	ctl.logger.Info("OAuth2SignInFacebookAction")

	auth := ctl.authConf.Facebook

	facebookOauthConfig.RedirectURL = auth.CallbackURL
	facebookOauthConfig.ClientID = auth.ClientID
	facebookOauthConfig.ClientSecret = auth.ClientSecret

	// token
	token := ctl.session.GenerateToken()
	ctl.session.SetToken(ctx, token)

	url := facebookOauthConfig.AuthCodeURL(token)

	// add display and auth_type
	// url = url + "&display=popup&auth_type=reauthenticate"
	url = url + "&display=popup"

	ctx.Redirect(http.StatusTemporaryRedirect, url) // 307
}

// OAuth2LoginAction is login by Google. (work in progress) [GET]
//func (ctl *controller) OAuth2LoginAction(ctx *gin.Context) {
//	ctl.logger.Info("OAuth2LoginAction")
//	/*
//		https://accounts.google.com/o/oauth2/auth?
//		scope=openid+email+profile&
//		state=G6OJI79YNaokmJNIGJRooGk4zUQVTRyi&
//		redirect_uri=https://courses.edx.org/auth/complete/google-oauth2/&
//		response_type=code&
//		client_id=370673641490-nd3s6q740l96uvk1vivsab65rltkgoc0.apps.googleusercontent.com
//	*/
//	//TODO:What is difference of parameter between sign in and login
//}

func checkError(ctx *gin.Context, logger *zap.Logger) bool {
	logger.Info("checkError")

	// When user choose access deny
	// http://localhost:9999/oauth2/callback?error=access_denied&state=66bc8679a5629423463943f679383b57
	if err, _ := ctx.GetQuery("error"); err != "" {
		logger.Error("query error", zap.Error(errors.New(err)))
		ctx.Redirect(http.StatusTemporaryRedirect, "/login") // 307
		return false
	}
	return true
}

func (ctl *controller) checkState(ctx *gin.Context, logger *zap.Logger) bool {
	logger.Info("checkState")

	state, _ := ctx.GetQuery("state")
	// ctl.logger.Debugf("state is %s", state)
	// ctl.logger.Debugf("saved state is %s", sess.GetTokenSession(ctx))
	if state == "" || ctl.session.GetToken(ctx) != state {
		// error
		logger.Error("checkState", zap.Error(errors.New("state is invalid")))
		ctx.Redirect(http.StatusTemporaryRedirect, "/") // 307
		return false
	}
	return true
}

func getToken(ctx *gin.Context, logger *zap.Logger, mode uint8) (token *oauth2.Token) {
	logger.Info("getToken")

	var err error

	code, _ := ctx.GetQuery("code")

	switch mode {
	case GoogleAuth:
		token, err = googleOauthConfig.Exchange(context.Background(), code)
	case FacebookAuth:
		token, err = facebookOauthConfig.Exchange(context.Background(), code)
	default:
		return nil
	}
	if err != nil {
		// error
		logger.Error("fail to call auth.Exchange", zap.Error(err))
		ctx.Redirect(http.StatusTemporaryRedirect, "/") // 307
		return nil
	}
	return token
}

func getUserInfo(ctx *gin.Context, logger *zap.Logger, token *oauth2.Token, res interface{}, mode uint8) bool {
	logger.Info("getUserInfo")

	var url string

	switch mode {
	case GoogleAuth:
		url = "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken
	case FacebookAuth:
		url = "https://graph.facebook.com/me?access_token=" + token.AccessToken
		url += "&fields=id,email,verified,name,first_name,last_name,link,picture,gender,locale"
		// client := facebookOauthConfig.Client(oauth2.NoContext, token)
		// response, err := client.Get("https://graph.facebook.com/me")
	default:
		return false
	}

	response, err := http.Get(url)
	if err != nil {
		// error
		logger.Error("getUserInfo: fail to call http.Get", zap.String("url", url), zap.Error(err))
		ctx.Redirect(http.StatusTemporaryRedirect, "/") // 307
		return false
	}

	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(contents, res)
	if err != nil {
		logger.Error("getUserInfo: fail to call json.Unmarshal", zap.Error(err))
		ctx.Redirect(http.StatusTemporaryRedirect, "/") // 307
		return false
	}

	return true
}

func (ctl *controller) registerOrLogin(ctx *gin.Context, mode uint8, user *user.User, userAuth *user.UserAuth) {
	ctl.logger.Info("registerOrLogin")

	if userAuth == nil {
		ctl.logger.Debug("registerOrLogin: no user")
		// 0:no user -> register and login
		id, err := ctl.userRepo.InsertUser(user)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		// Session
		ctl.session.SetUserID(ctx, int(id))

	} else {
		ctl.logger.Debug("registerOrLogin", zap.Any("user", userAuth))
		// oauth_flg //0, 1:google, 2:facebook
		if userAuth.ID != 0 && userAuth.Auth == mode {
			// 1:existing user (google) -> login procedure
			// Session
			ctl.session.SetUserID(ctx, userAuth.ID)
		} else {
			ctl.logger.Debug("registerOrLogin: redirect login page. user is already existing")
			// 2:existing user (no auth or another auth) -> err
			ctx.Redirect(http.StatusTemporaryRedirect, "/login") // 307
			return
		}
	}

	// Login
	// token delete
	ctl.session.DeleteToken(ctx)

	// Redirect[GET]
	ctx.Redirect(http.StatusTemporaryRedirect, "/accounts") // 307
}

// OAuth2CallbackGoogleAction is callback from Google[GET]
func (ctl *controller) OAuth2CallbackGoogleAction(ctx *gin.Context) {
	ctl.logger.Info("OAuth2CallbackGoogleAction")
	mode := GoogleAuth

	// 0.check deny
	bRet := checkError(ctx, ctl.logger)
	if !bRet {
		return
	}

	// 1.Confirm State(token)
	bRet = ctl.checkState(ctx, ctl.logger)
	if !bRet {
		return
	}

	// 2.connection server to server
	token := getToken(ctx, ctl.logger, mode)
	if token == nil {
		return
	}

	// 3.get user info
	resGoogle := ResGoogle{}
	bRet = getUserInfo(ctx, ctl.logger, token, &resGoogle, mode)
	if !bRet {
		return
	}

	ctl.logger.Debug("OAuth2CallbackGoogleAction", zap.Any("response body", resGoogle))

	// 4.check Email
	ctl.logger.Debug("", zap.String("email", resGoogle.Email))
	userAuth, err := ctl.userRepo.OAuth2Login(resGoogle.Email)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	// 5. register or login
	user := &user.User{
		FirstName: resGoogle.FirstName,
		LastName:  resGoogle.LastName,
		Email:     resGoogle.Email,
		Password:  "google-password",
		OAuth2:    mode,
	}
	ctl.registerOrLogin(ctx, mode, user, userAuth)
}

// OAuth2CallbackFacebookAction is callback from Facebook [GET]
func (ctl *controller) OAuth2CallbackFacebookAction(ctx *gin.Context) {
	ctl.logger.Info("OAuth2CallbackFacebookAction")
	mode := FacebookAuth

	// 0.check deny
	bRet := checkError(ctx, ctl.logger)
	if !bRet {
		return
	}

	// 1.Confirm State(token)
	bRet = ctl.checkState(ctx, ctl.logger)
	if !bRet {
		return
	}

	// 2.connection server to server
	token := getToken(ctx, ctl.logger, mode)
	if token == nil {
		return
	}

	// 3.get user info
	resFacebook := ResFacebook{}
	bRet = getUserInfo(ctx, ctl.logger, token, &resFacebook, mode)
	if !bRet {
		return
	}

	ctl.logger.Debug("OAuth2CallbackFacebookAction", zap.Any("response body", resFacebook))
	// img := "https://graph.facebook.com/" + id + "/picture?width=180&height=180"

	// 4.check Email
	ctl.logger.Debug("OAuth2CallbackFacebookAction", zap.String("email", resFacebook.Email))
	userAuth, err := ctl.userRepo.OAuth2Login(resFacebook.Email)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	// 5. register or login
	user := &user.User{
		FirstName: resFacebook.FirstName,
		LastName:  resFacebook.LastName,
		Email:     resFacebook.Email,
		Password:  "facebook-password",
		OAuth2:    mode,
	}
	ctl.registerOrLogin(ctx, mode, user, userAuth)
}
