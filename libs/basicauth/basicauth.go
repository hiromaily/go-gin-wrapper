package basicauth

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const AuthUserKey = "user"

type (
	Accounts map[string]string
	authPair struct {
		Value string
		User  string
	}
	authPairs []authPair
)

func (a authPairs) searchCredential(authValue string) (string, bool) {
	if len(authValue) == 0 {
		return "", false
	}
	for _, pair := range a {
		if pair.Value == authValue {
			return pair.User, true
		}
	}
	return "", false
}

// BasicAuthForRealm returns a Basic HTTP Authorization middleware. It takes as arguments a map[string]string where
// the key is the user name and the value is the password, as well as the name of the Realm.
// If the realm is empty, "Authorization Required" will be used by default.
// (see http://tools.ietf.org/html/rfc2617#section-1.2)
func BasicAuthForRealm(accounts Accounts, realm string) gin.HandlerFunc {
	if realm == "" {
		realm = "Authorization Required"
	}
	realm = "Basic realm=" + strconv.Quote(realm)
	pairs := processAccounts(accounts)
	return func(c *gin.Context) {
		// Search user in the slice of allowed credentials
		user, found := pairs.searchCredential(c.Request.Header.Get("Authorization"))
		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(401)
			//c.Status(code)
			//c.Abort()
			//TODO:return specific error html
			fmt.Println("[BasicAuthForRealm] before return html")
			c.HTML(http.StatusUnauthorized, "pages/errors/error.tmpl", gin.H{
				"message": "401 errors",
			})

		} else {
			// The user credentials was found, set user's id to key AuthUserKey in this context, the userId can be read later using
			// c.MustGet(gin.AuthUserKey)
			c.Set(AuthUserKey, user)
		}
	}
}

// BasicAuth returns a Basic HTTP Authorization middleware. It takes as argument a map[string]string where
// the key is the user name and the value is the password.
func BasicAuth(accounts Accounts) gin.HandlerFunc {
	return BasicAuthForRealm(accounts, "")
}

func processAccounts(accounts Accounts) authPairs {
	assert1(len(accounts) > 0, "Empty list of authorized credentials")
	pairs := make(authPairs, 0, len(accounts))
	for user, password := range accounts {
		assert1(len(user) > 0, "User can not be empty")
		value := authorizationHeader(user, password)
		pairs = append(pairs, authPair{
			Value: value,
			User:  user,
		})
	}
	return pairs
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}

func secureCompare(given, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	}
	/* Securely compare actual to itself to keep constant time, but always return false */
	return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}
