package regexp

import (
	"regexp"
)

func CheckRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}

func IsInvisiblefile(target string) bool {
	return CheckRegexp(`^[\\.].*$`, target)
}

func IsGoFile(target string) bool {
	return CheckRegexp(`^.*\.go$`, target)
}

func IsTmplFile(target string) bool {
	return CheckRegexp(`^.*\.tmpl$`, target)
}

func IsExtFile(target, ext string) bool {
	return CheckRegexp(`^.*\.`+ext+`$`, target)
}

func IsHeaderURL(target string) bool {
	return CheckRegexp(`^http(s)?:\/\/`, target)
}

func IsBenchTest(target string) bool {
	return CheckRegexp(`^-test.bench`, target)
}
