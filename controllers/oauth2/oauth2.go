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
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

type ResGoogle struct {
	Id            string `json:"id"`
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
)

const GoogleAuth string = "1"

//Sign In [GET]
func SignInAction(c *gin.Context) {
	lg.Info("SignInAction()")

	auth := conf.GetConf().Auth.Google

	googleOauthConfig.RedirectURL = auth.CallbackURL
	googleOauthConfig.ClientID = auth.ClientID
	googleOauthConfig.ClientSecret = auth.ClientSecret

	//token
	token := csrf.CreateToken()
	sess.SetTokenSession(c, token)

	url := googleOauthConfig.AuthCodeURL(token)
	//http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	c.Redirect(http.StatusTemporaryRedirect, url) //307
}

//Login [GET]
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

//Callback [GET]
func CallbackAction(c *gin.Context) {
	lg.Info("CallbackAction()")

	//1.Confirm State(token)
	state, _ := c.GetQuery("state")
	//lg.Debugf("state is %s", state)
	//lg.Debugf("saved state is %s", sess.GetTokenSession(c))
	if state == "" || sess.GetTokenSession(c) != state {
		//error
		lg.Error("state is invalid.")
		c.Redirect(http.StatusTemporaryRedirect, "/") //307
		return
	}

	//2.connection server to server
	code, _ := c.GetQuery("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		//error
		lg.Errorf("Code exchange failed with '%s'", err)
		c.Redirect(http.StatusTemporaryRedirect, "/") //307
		return
	}

	//3.get user info
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		//error
		lg.Errorf("Get user info failed with '%s'", err)
		c.Redirect(http.StatusTemporaryRedirect, "/") //307
		return
	}

	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)

	resGoogle := ResGoogle{}
	err = json.Unmarshal(contents, &resGoogle)
	if err != nil {
		lg.Errorf("Parse of response as json failed with '%s'", err)
		c.Redirect(http.StatusTemporaryRedirect, "/") //307
		return
	}

	lg.Debugf("response body is %+s", resGoogle)

	//4.check Email
	lg.Debugf("email is %s", resGoogle.Email)
	userAuth, err := models.GetDB().OauthLogin(resGoogle.Email)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if userAuth == nil {
		lg.Debug("no user on t_users")
		//0:no user -> register and login
		user := &models.Users{
			FirstName: resGoogle.FirstName,
			LastName:  resGoogle.LastName,
			Email:     resGoogle.Email,
			Password:  "google-password",
			Oauth2Flg: GoogleAuth,
		}

		lg.Debug("InsertUser()")
		id, err := models.GetDB().InsertUser(user)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		//Session
		sess.SetUserSession(c, int(id))

	} else {
		lg.Debug("There is user: %v", userAuth)
		//oauth_flg //0, 1:google, 2:facebook
		if userAuth.Id != 0 && userAuth.Auth == GoogleAuth {
			lg.Debug("login proceduer")
			//1:existing user (google) -> login procedure
			//Session
			sess.SetUserSession(c, userAuth.Id)
		} else {
			lg.Debug("redirect login page")
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
