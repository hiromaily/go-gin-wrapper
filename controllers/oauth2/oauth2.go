package oauth2

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	"github.com/hiromaily/go-gin-wrapper/libs/csrf"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
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

	//TODO: 成功時は登録成功ページにジャンプ
	//View
	c.HTML(http.StatusOK, "pages/oauth2/callback.tmpl", gin.H{
		"title": "News Page",
	})

}
