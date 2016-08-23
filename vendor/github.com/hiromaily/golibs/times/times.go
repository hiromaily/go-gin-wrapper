package times

import (
	"fmt"
	"strings"
	"time"
)

//format
const (
	// 日付のフォーマット
	FORMAT_A string = "1/2"
	FORMAT_B string = "【1/2】"
	FORMAT_C string = "【1月2日】"

	// 日付のフォーマット(+曜日)
	FORMAT_A_WEEK string = "1/2(%s)"
	FORMAT_B_WEEK string = "【1/2(%s)】"
	FORMAT_C_WEEK string = "【1月2日(%s)】"
)

// day of week
var JAPANESE_WEEKDAYS = []string{"日", "月", "火", "水", "木", "金", "土"}

var TimeLayouts = []string{
	"Mon, _2 Jan 2006 15:04:05 MST",   //0
	"Mon, _2 Jan 2006 15:04:05 -0700", //1
	time.ANSIC,                        //2
	time.UnixDate,                     //3
	time.RubyDate,                     //4
	time.RFC822,                       //5
	time.RFC822Z,                      //6
	time.RFC850,                       //7
	time.RFC1123,                      //8
	time.RFC1123Z,                     //9
	time.RFC3339,                      //10
	time.RFC3339Nano,                  //11
}

//[1 2 3 4 5 6 7 9 10 11]
//[0 2 3 4 5 6 7 8 10 11]

//return accessible format.
func CheckParseTime(s string) []int {
	s = strings.TrimSpace(s)
	iRet := []int{}

	for i, layout := range TimeLayouts {
		_, err := time.Parse(layout, s)
		if err == nil {
			iRet = append(iRet, i)
		}
	}

	return iRet
}

func ParseTime(str string) (t time.Time, err error) {
	str = strings.TrimSpace(str)

	for _, layout := range TimeLayouts {
		t, err = time.Parse(layout, str)
		if err == nil {
			return t, err
		}
	}

	return time.Time{}, err
}

//Formatter Date
func GetFormatDate(strDate string, format string, addWeek bool) string {
	//2016-05-13 16:52:49
	//To time object
	t, _ := time.Parse("2006-01-02 15:04:05", strDate)

	//t.Month()
	//t.Day()

	//Format
	var baseFormat string
	if addWeek {
		if format != "" {
			baseFormat = fmt.Sprintf(format, JAPANESE_WEEKDAYS[t.Weekday()])
		} else {
			baseFormat = fmt.Sprintf(FORMAT_A_WEEK, JAPANESE_WEEKDAYS[t.Weekday()])
		}
	} else {
		if format != "" {
			baseFormat = format
		} else {
			baseFormat = FORMAT_A
		}
	}

	return t.Format(baseFormat)
}

//Formatter Time
func GetFormatTime(strTime string, format string) string {
	t, _ := time.Parse("2006-01-02 15:04:05", strTime)

	//t.Hour()
	//t.Minute()
	//t.Second()

	if format != "" {
		return t.Format(format)
	} else {
		return t.Format("15:04")
	}

}

func GetCurrentTimeByStr() string {
	t := time.Now()
	layout := "2006-01-02 15:04:05"
	return t.Format(layout)
}

func PerseTimeForLastModified(lastModified string) (time.Time, error) {
	//Tue, 16 Aug 2016 01:31:09 GMT
	return time.Parse(time.RFC1123, lastModified)
}

func PerseTimeForRss(str string) (time.Time, error) {
	t, err := time.Parse(time.RFC1123, str)
	if err != nil {
		t, err = time.Parse(time.RFC1123Z, str)
	}
	return t, err
}

//Timer
//Caller: defer Track(time.Now(), "parseFile()")
//https://medium.com/@2xb/execution-time-tracking-in-golang-9379aebfe20e#.ffxgxejim
func Track(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
