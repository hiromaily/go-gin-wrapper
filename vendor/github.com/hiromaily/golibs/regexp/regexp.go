package regexp

import (
	"regexp"
)

// CheckRegexp is check str using pattern reg
func CheckRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

// IsInvisiblefile is whether target is invisible file or not
func IsInvisiblefile(target string) bool {
	return CheckRegexp(`^[\\.].*$`, target)
}

// IsGoFile is .go or not
func IsGoFile(target string) bool {
	return CheckRegexp(`^.*\.go$`, target)
}

// IsTmplFile is .tmpl or not
func IsTmplFile(target string) bool {
	return CheckRegexp(`^.*\.tmpl$`, target)
}

// IsStaticFile is static file which has extension in file name or not
func IsStaticFile(target string) bool {
	//is there any suffix
	//.+\.(csv|pdf)
	return CheckRegexp(`^.*\.`, target)
}

// IsExtFile is check that target include ext string
func IsExtFile(target, ext string) bool {
	return CheckRegexp(`^.*\.`+ext+`$`, target)
}

// IsHeaderURL is check url
func IsHeaderURL(target string) bool {
	return CheckRegexp(`^http(s)?:\/\/`, target)
}

// IsBenchTest is check parameter for bench test
func IsBenchTest(target string) bool {
	return CheckRegexp(`^-test.bench`, target)
}
