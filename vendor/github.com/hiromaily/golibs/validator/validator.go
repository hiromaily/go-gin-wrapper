package validator

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// BasicValidator is validator function
type BasicValidator func(str string) bool

// CalcValidator is validator function
type CalcValidator func(str string, num int) bool

// TagMap is to map tag name of struct to function
var TagMap = map[string]BasicValidator{
	"nonempty": isNonEmpty,
	"email":    isEmail,
	"url":      isURL,
	"number":   isNumber,
	"alphabet": isAlphabet,
}

// TagMapCal is to map tag name of struct to function
var TagMapCal = map[string]CalcValidator{
	"min": isMinOK,
	"max": isMaxOK,
}

//-----------------------------------------------------------------------------
// functions for validator
// github.com/asaskevich/govalidator/validator.go
//-----------------------------------------------------------------------------

func isNonEmpty(str string) bool {
	return str != ""
}

func isEmail(str string) bool {
	regEx := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regEx.MatchString(str)
}

func isURL(str string) bool {
	if str == "" || len(str) >= 2083 || len(str) <= 3 || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	regEx := regexp.MustCompile(`http(s)?://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?`)
	return regEx.MatchString(str)
}

func isNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func isAlphabet(str string) bool {
	for i := range str {
		if str[i] < 'A' || str[i] > 'z' {
			return false
		} else if str[i] > 'Z' && str[i] < 'a' {
			return false
		}
	}
	return true
}

func isMinOK(str string, num int) bool {
	return utf8.RuneCountInString(str) >= num
}

func isMaxOK(str string, num int) bool {
	return utf8.RuneCountInString(str) <= num
}

func getErrorMsgFmt(chkItem string, errFmt map[string]string) (string, error) {
	if _, ok := errFmt[chkItem]; ok {
		return errFmt[chkItem], nil
	}
	return "", errors.New("Not found key")
}

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------

func checkValidation(str string, val string, disp string) []string {
	ret := []string{disp}
	strs := strings.Split(str, ",")
	for _, v := range strs {
		//When included「=」on v, divide it.
		equals := strings.Split(v, "=")
		var bRet bool
		if len(equals) > 1 {
			num, _ := strconv.Atoi(equals[1])
			bRet = TagMapCal[equals[0]](val, num)
			if !bRet {
				ret = append(ret, v)
			}
		} else {
			//nonempty
			bRet = TagMap[v](val)
			if !bRet {
				ret = append(ret, v)
			}
		}
	}
	//[]string{"min"}
	return ret
}

// CheckValidation to check validation after extracted tag from struct type.
func CheckValidation(s interface{}, brankSkip bool) map[string][]string {

	mRet := make(map[string][]string)

	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {

		field := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		valid := tag.Get("valid")
		fld := tag.Get("field")
		disp := tag.Get("dispName")

		if valid != "" {
			//When check is required, check specific field
			val, _ := field.Interface().(string)
			if !brankSkip || (brankSkip && val != "") {
				//Returned value is slice, stored name of error
				ret := checkValidation(valid, val, disp)
				if len(ret) > 1 {
					mRet[fld] = ret
				}
			}
		}
	}

	return mRet
}

// ConvertErrorMsgs is to convert error messages
func ConvertErrorMsgs(data map[string][]string, errFmt map[string]string) []string {

	msgs := []string{}
	for key, val := range data {
		fmtkey := strings.Split(val[1], "=")
		strFmt, err := getErrorMsgFmt(fmtkey[0], errFmt)
		if err == nil {
			if len(fmtkey) == 1 {
				msgs = append(msgs, fmt.Sprintf(strFmt, data[key][0]))
			} else {
				msgs = append(msgs, fmt.Sprintf(strFmt, fmtkey[1], data[key][0]))
			}
		}
		//In case of error, skip.
	}
	return msgs
}
