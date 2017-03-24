package login

import (
	"errors"
	models "github.com/hiromaily/go-gin-wrapper/models/mysql"
	lg "github.com/hiromaily/golibs/log"
	valid "github.com/hiromaily/golibs/validator"
	//gin "gopkg.in/gin-gonic/gin.v1"
	"github.com/gin-gonic/gin"
)

// Request is request structure for login
type Request struct {
	Email string `valid:"nonempty,email,min=5,max=40" field:"email" dispName:"E-Mail"`
	Pass  string `valid:"nonempty,min=8,max=20" field:"pass" dispName:"Password"`
}

// ErrFmt is error message for validation package
// TODO: Should this be defined as something common library?
var ErrFmt = map[string]string{
	"nonempty": "Empty is not allowed on %s",
	"email":    "Format of %s is invalid",
	"alphanum": "Only alphabet is allowd on %s",
	"min":      "At least %s of characters is required on %s",
	"max":      "At a maximum %s of characters is allowed on %s",
}

// CheckLoginOnHTML is check login on html page.
func CheckLoginOnHTML(c *gin.Context) (int, *Request, []string) {
	//Get Post Parameters
	posted := &Request{
		Email: c.PostForm("inputEmail"),
		Pass:  c.PostForm("inputPassword"),
	}

	//Validation
	mRet := valid.CheckValidation(posted, false)
	if len(mRet) != 0 {
		errs := valid.ConvertErrorMsgs(mRet, ErrFmt)
		lg.Debugf("validation error: %#v", errs)

		//return
		//resLogin(c, posted, "", errs)
		return 0, posted, errs
	}

	//Check inputed mail and password
	userID, err := models.GetDB().IsUserEmail(posted.Email, posted.Pass)
	if err != nil {
		errs := []string{"E-mail or Password is made a mistake."}
		lg.Debugf("login error: %#v", errs)

		//return
		//resLogin(c, posted, "", errs)
		return 0, posted, errs
	}
	return userID, nil, nil
}

// CheckLoginOnAPI is check login on API
func CheckLoginOnAPI(c *gin.Context) (int, string, error) {
	posted := &Request{
		Email: c.PostForm("inputEmail"),
		Pass:  c.PostForm("inputPassword"),
	}

	//Validation
	mRet := valid.CheckValidation(posted, false)
	if len(mRet) != 0 {
		return 0, "", errors.New("validation error")
	}

	//Check inputed mail and password
	userID, err := models.GetDB().IsUserEmail(posted.Email, posted.Pass)
	if err != nil {
		return 0, "", errors.New("login error")
	}
	return userID, c.PostForm("inputEmail"), nil
}
