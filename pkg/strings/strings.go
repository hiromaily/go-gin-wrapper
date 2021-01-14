package strings

import (
	"strconv"
	"strings"
)

// SearchIndex is to search string
func SearchIndex(target string, sources []string) int {
	retIdx := -1
	if len(sources) == 0 {
		return retIdx
	}
	for i, val := range sources {
		if val == target {
			retIdx = i
			break
		}
	}
	return retIdx
}

// SearchIndexLower doesn't distinguish  upper case and lower case
func SearchIndexLower(target string, sources []string) int {
	retIdx := -1
	if len(sources) == 0 {
		return retIdx
	}
	for i, val := range sources {
		if strings.EqualFold(val, target) {
			retIdx = i
			break
		}
	}
	return retIdx
}

func Itos(val interface{}) string {
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}

func Atoi(str string) (ret int) {
	ret, _ = strconv.Atoi(str)
	return
}
