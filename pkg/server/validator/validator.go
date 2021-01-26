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

// refer to `github.com/asaskevich/govalidator/validator.go`

//-----------------------------------------------------------------------------
// BasicValidator is validator function
//-----------------------------------------------------------------------------

// BasicValidator is basic validator
type BasicValidator func(str string) bool

// TagMap is to map tag name of struct to function
var TagMap = map[string]BasicValidator{
	"nonempty": isNonEmpty,
	"email":    isEmail,
	"url":      isURL,
	"number":   isNumber,
	"alphabet": isAlphabet,
}

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

//-----------------------------------------------------------------------------
// CalcValidator is validator function
//-----------------------------------------------------------------------------

// CalcValidator is validator with parameter
type CalcValidator func(str string, num int) bool

// TagMapCal is to map tag name of struct to function
var TagMapCal = map[string]CalcValidator{
	"min": isMinOK,
	"max": isMaxOK,
}

func isMinOK(str string, num int) bool {
	return utf8.RuneCountInString(str) >= num
}

func isMaxOK(str string, num int) bool {
	return utf8.RuneCountInString(str) <= num
}

func getErrorMsgFmt(chkItem string, errFormat map[string]string) (string, error) {
	if _, ok := errFormat[chkItem]; ok {
		return errFormat[chkItem], nil
	}
	return "", errors.New("not found key")
}

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------

// Validate validates structure with tag
func Validate(s interface{}, brankSkip bool) map[string][]string {
	failedFields := make(map[string][]string)

	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {

		field := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		// tag
		if validAll := tag.Get("valid"); validAll != "" {
			// tag
			fld := tag.Get("field")
			disp := tag.Get("dispName")

			// check specific field when required,
			val, _ := field.Interface().(string)
			if !brankSkip || (brankSkip && val != "") {
				if invalid := validate(validAll, val, disp); len(invalid) > 1 {
					failedFields[fld] = invalid
				}
			}
		}
	}
	return failedFields
}

func validate(validAll, val, disp string) []string {
	failed := []string{disp}
	valid := strings.Split(validAll, ",")
	for _, v := range valid {
		// when `v` included`=`, divide it
		equals := strings.Split(v, "=")
		if len(equals) > 1 {
			num, _ := strconv.Atoi(equals[1])
			if ok := TagMapCal[equals[0]](val, num); !ok {
				failed = append(failed, v)
			}
		} else {
			// nonempty
			if ok := TagMap[v](val); !ok {
				failed = append(failed, v)
			}
		}
	}
	return failed
}

// ConvertErrorMsgs converts error messages
func ConvertErrorMsgs(data map[string][]string, errFormat map[string]string) []string {
	msgs := make([]string, 0)
	for key, val := range data {
		fmtkey := strings.Split(val[1], "=")
		strFmt, err := getErrorMsgFmt(fmtkey[0], errFormat)
		if err == nil {
			if len(fmtkey) == 1 {
				msgs = append(msgs, fmt.Sprintf(strFmt, data[key][0]))
			} else {
				msgs = append(msgs, fmt.Sprintf(strFmt, fmtkey[1], data[key][0]))
			}
		}
	}
	return msgs
}
