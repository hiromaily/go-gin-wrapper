package time

import (
	"fmt"
	"strings"
	"time"
)

// JapaneseWeedDays is day of week
var JapaneseWeedDays = []string{"日", "月", "火", "水", "木", "金", "土"}

// TimeLayouts is time layouts
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

// CheckParseTime is to return accessible format
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

// ParseTime is to parse time by available time format.
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

// ParseTimeForLastModified is to parse time for LastModified
func ParseTimeForLastModified(lastModified string) (time.Time, error) {
	//Tue, 16 Aug 2016 01:31:09 GMT
	return time.Parse(time.RFC1123, lastModified)
}

// ParseTimeForRss is to parse time for RSS
func ParseTimeForRss(str string) (time.Time, error) {
	t, err := time.Parse(time.RFC1123, str)
	if err != nil {
		t, err = time.Parse(time.RFC1123Z, str)
	}
	return t, err
}

// Track is to track elapsed time
//  e.g. Caller: defer Track(time.Now(), "parseFile()")
//  https://medium.com/@2xb/execution-time-tracking-in-golang-9379aebfe20e#.ffxgxejim
func Track(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

// GetCurrentDateTimeByStr is to get current time by string
func GetCurrentDateTimeByStr(format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	t := time.Now()
	return t.Format(format)
}

// GetFormatDate is to get format date
func GetFormatDate(strDate string, format string, addWeek bool) string {

	t, _ := time.Parse("2006-01-02 15:04:05", strDate)
	//t.Month()
	//t.Day()

	//Format
	if format == "" && !addWeek {
		format = "【1/2】"
	} else if format == "" && addWeek {
		format = fmt.Sprintf("【1/2(%s)】", JapaneseWeedDays[t.Weekday()])
	} else if format != "" && addWeek {
		format = fmt.Sprintf(format, JapaneseWeedDays[t.Weekday()])
	}

	return t.Format(format)
}

// GetFormatTime is to format time
func GetFormatTime(strTime string, format string) string {
	if format == "" {
		format = "15:04"
	}

	t, _ := time.Parse("2006-01-02 15:04:05", strTime)
	//t.Hour()
	//t.Minute()
	//t.Second()

	return t.Format(format)
}
