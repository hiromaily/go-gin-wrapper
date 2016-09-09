package oauth2

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	"github.com/hiromaily/go-gin-wrapper/libs/csrf"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	models "github.com/hiromaily/go-gin-wrapper/models/mysql"
	lg "github.com/hiromaily/golibs/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

// ResGoogle is for response data from google
type ResGoogle struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`        //full name
	FirstName     string `json:"given_name"`  //first name
	LastName      string `json:"family_name"` //last name
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
	Name          string `json:"name"`       //full name
	FirstName     string `json:"first_name"` //first name
	LastName      string `json:"last_name"`  //last name
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
	GoogleAuth string = "1"
	// FacebookAuth is for Facebook
	FacebookAuth string = "2"
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
		//Scopes:   []string{"first_name", "last_name", "link", "picture", "email"},
		Endpoint: facebook.Endpoint,
	}
)

// SignInGoogleAction is sign in by Google [GET]
func SignInGoogleAction(c *gin.Context) {
	lg.Info("SignInGoogleAction()")

	auth := conf.GetConf().Auth.Google

	googleOauthConfig.RedirectURL = auth.CallbackURL
	googleOauthConfig.ClientID = auth.ClientID
	googleOauthConfig.ClientSecret = auth.ClientSecret

	//token
	token := csrf.CreateToken()
	sess.SetTokenSession(c, token)

	url := googleOauthConfig.AuthCodeURL(token)
	c.Redirect(http.StatusTemporaryRedirect, url) //307
}

// SignInFacebookAction is sign in by Facebook [GET]
func SignInFacebookAction(c *gin.Context) {
	lg.Info("SignInFacebookAction()")

	auth := conf.GetConf().Auth.Facebook

	facebookOauthConfig.RedirectURL = auth.CallbackURL
	facebookOauthConfig.ClientID = auth.ClientID
	facebookOauthConfig.ClientSecret = auth.ClientSecret

	//token
	token := csrf.CreateToken()
	sess.SetTokenSession(c, token)

	url := facebookOauthConfig.AuthCodeURL(token)

	//add display and auth_type
	//url = url + "&display=popup&auth_type=reauthenticate"
	url = url + "&display=popup"

	c.Redirect(http.StatusTemporaryRedirect, url) //307
}

// LoginAction is login by Google. (work in progress) [GET]
func LoginAction(c *gin.Context) {
	lg.Info("LoginAction()")
	/*
		https://accounts.google.com/o/oauth2/auth?
		scope=openid+email+profile&
		state=G6OJI79YNaokmJNIGJRooGk4zUQVTRyi&
		redirect_uri=https://courses.edx.org/auth/complete/google-oauth2/&
		response_type=code&
		client_id=370673641490-nd3s6q740l96uvk1vivsab65rltkgoc0.apps.googleusercontent.com
	*/
	//TODO:What is difference of parameter between sign in and login
}

func checkError(c *gin.Context) bool {
	lg.Info("checkError()")

	//When user choose access deny
	//http://localhost:9999/oauth2/callback?error=access_denied&state=66bc8679a5629423463943f679383b57
	qeyErr, _ := c.GetQuery("error")
	if qeyErr != "" {
		lg.Debugf("error is %s", qeyErr)
		c.Redirect(http.StatusTemporaryRedirect, "/login") //307
		return false
	}
	return true
}

func checkState(c *gin.Context) bool {
	lg.Info("checkState()")

	state, _ := c.GetQuery("state")
	//lg.Debugf("state is %s", state)
	//lg.Debugf("saved state is %s", sess.GetTokenSession(c))
	if state == "" || sess.GetTokenSession(c) != state {
		//error
		lg.Error("state is invalid.")
		c.Redirect(http.StatusTemporaryRedirect, "/") //307
		return false
	}
	return true
}

func getToken(c *gin.Context, mode string) (token *oauth2.Token) {
	lg.Info("getToken()")

	var err error

	code, _ := c.GetQuery("code")

	switch mode {
	case GoogleAuth:
		token, err = googleOauthConfig.Exchange(oauth2.NoContext, code)
	case FacebookAuth:
		token, err = facebookOauthConfig.Exchange(oauth2.NoContext, code)
	default:
		return nil
	}

	if err != nil {
		//error
		lg.Errorf("Code exchange failed with '%s'", err)
		c.Redirect(http.StatusTemporaryRedirect, "/") //307
		return nil
	}
	return token
}

func getUserInfo(c *gin.Context, token *oauth2.Token, res interface{}, mode string) bool {
	lg.Info("getUserInfo()")

	var url string

	switch mode {
	case GoogleAuth:
		url = "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken
	case FacebookAuth:
		url = "https://graph.facebook.com/me?access_token=" + token.AccessToken
		url += "&fields=id,email,verified,name,first_name,last_name,link,picture,gender,locale"
		//client := facebookOauthConfig.Client(oauth2.NoContext, token)
		//response, err := client.Get("https://graph.facebook.com/me")
	default:
		return false
	}

	response, err := http.Get(url)
	if err != nil {
		//error
		lg.Errorf("Get user info failed with '%s'", err)
		c.Redirect(http.StatusTemporaryRedirect, "/") //307
		return false
	}

	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(contents, res)
	if err != nil {
		lg.Errorf("Parse of response as json failed with '%s'", err)
		c.Redirect(http.StatusTemporaryRedirect, "/") //307
		return false
	}

	return true
}

func registerOrLogin(c *gin.Context, mode string, uA *models.UserAuth, user *models.Users) {
	lg.Info("registerOrLogin()")

	if uA == nil {
		lg.Debug("no user on t_users")
		//0:no user -> register and login

		lg.Debug("InsertUser()")
		id, err := models.GetDB().InsertUser(user)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		//Session
		sess.SetUserSession(c, int(id))

	} else {
		lg.Debug("There is user: %v", uA)
		//oauth_flg //0, 1:google, 2:facebook
		if uA.ID != 0 && uA.Auth == mode {
			lg.Debug("login proceduer")
			//1:existing user (google) -> login procedure
			//Session
			sess.SetUserSession(c, uA.ID)
		} else {
			lg.Debug("redirect login page. user is already exsisting.")
			//2:existing user (no auth or another auth) -> err
			c.Redirect(http.StatusTemporaryRedirect, "/login") //307
			return
		}
	}

	//Login
	//token delete
	sess.DelTokenSession(c)

	//Redirect[GET]
	c.Redirect(http.StatusTemporaryRedirect, "/accounts") //307

	return
}

// CallbackGoogleAction is callback from Google[GET]
func CallbackGoogleAction(c *gin.Context) {
	lg.Info("CallbackGoogleAction()")
	mode := GoogleAuth

	//0.check deny
	bRet := checkError(c)
	if !bRet {
		return
	}

	//1.Confirm State(token)
	bRet = checkState(c)
	if !bRet {
		return
	}

	//2.connection server to server
	token := getToken(c, mode)
	if token == nil {
		return
	}

	//3.get user info
	resGoogle := ResGoogle{}
	bRet = getUserInfo(c, token, &resGoogle, mode)
	if !bRet {
		return
	}

	lg.Debugf("response body is %+s", resGoogle)

	//4.check Email
	lg.Debugf("email is %s", resGoogle.Email)
	userAuth, err := models.GetDB().OAuth2Login(resGoogle.Email)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	//5. register or login
	user := &models.Users{
		FirstName: resGoogle.FirstName,
		LastName:  resGoogle.LastName,
		Email:     resGoogle.Email,
		Password:  "google-password",
		OAuth2Flg: mode,
	}
	registerOrLogin(c, mode, userAuth, user)

	return
}

// CallbackFacebookAction is callback from Facebook [GET]
func CallbackFacebookAction(c *gin.Context) {
	lg.Info("CallbackFacebookAction()")
	mode := FacebookAuth

	//0.check deny
	bRet := checkError(c)
	if !bRet {
		return
	}

	//1.Confirm State(token)
	bRet = checkState(c)
	if !bRet {
		return
	}

	//2.connection server to server
	token := getToken(c, mode)
	if token == nil {
		return
	}

	//3.get user info
	resFacebook := ResFacebook{}
	bRet = getUserInfo(c, token, &resFacebook, mode)
	if !bRet {
		return
	}

	lg.Debugf("response body is %+s", resFacebook)
	//img := "https://graph.facebook.com/" + id + "/picture?width=180&height=180"

	//4.check Email
	lg.Debugf("email is %s", resFacebook.Email)
	userAuth, err := models.GetDB().OAuth2Login(resFacebook.Email)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	//5. register or login
	user := &models.Users{
		FirstName: resFacebook.FirstName,
		LastName:  resFacebook.LastName,
		Email:     resFacebook.Email,
		Password:  "facebook-password",
		OAuth2Flg: mode,
	}
	registerOrLogin(c, mode, userAuth, user)

	return
}
