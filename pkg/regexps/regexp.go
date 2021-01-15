package regexps

import (
	"regexp"
)

// IsStaticFile is static file which has extension in file name or not
func IsStaticFile(target string) bool {
	//is there any suffix
	//.+\.(csv|pdf)
	return checkRegexp(`^.*\.`, target)
}

// IsInvisiblefile is whether target is invisible file or not
func IsInvisiblefile(target string) bool {
	return checkRegexp(`^[\\.].*$`, target)
}

// IsExtFile is check that target include ext string
func IsExtFile(target, ext string) bool {
	return checkRegexp(`^.*\.`+ext+`$`, target)
}

// checkRegexp is check str using pattern reg
func checkRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}
