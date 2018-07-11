package regexp

import (
	re "regexp"
)

// CheckRegexp is check str using pattern reg
func CheckRegexp(reg, str string) bool {
	return re.MustCompile(reg).Match([]byte(str))
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

func Replace(path, from, to string) string {
	//reg := re.MustCompile(from)
	//return reg.ReplaceAllString(path, to)
	reg, _ := re.Compile(from)
	if reg.MatchString(path) {
		return reg.ReplaceAllString(path, to)
	}
	return "error"
}

func Replace2(path, from, to string) string {
	reg := re.MustCompile(from)
	return reg.ReplaceAllString(path, to)
}
