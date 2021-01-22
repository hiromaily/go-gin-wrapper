package strings

import (
	"strconv"
	"strings"
)

// SearchIndex searches target string from list
func SearchIndex(target string, sources []string) int {
	idx := -1
	if len(sources) == 0 {
		return idx
	}
	for i, val := range sources {
		if val == target {
			idx = i
			break
		}
	}
	return idx
}

// SearchIndexLower searches target string from list
// which doesn't distinguish  upper case and lower case
func SearchIndexLower(target string, sources []string) int {
	idx := -1
	if len(sources) == 0 {
		return idx
	}
	for i, source := range sources {
		if strings.EqualFold(source, target) {
			idx = i
			break
		}
	}
	return idx
}

// Itos converts interface to string
func Itos(target interface{}) string {
	str, ok := target.(string)
	if !ok {
		return ""
	}
	return str
}

// Atoi converts string to int
// take care to use this func because error is ignored
func Atoi(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}
