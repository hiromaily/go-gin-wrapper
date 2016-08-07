package times

import (
	"fmt"
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

//Timer
//Caller: defer Track(time.Now(), "parseFile()")
//https://medium.com/@2xb/execution-time-tracking-in-golang-9379aebfe20e#.ffxgxejim
func Track(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
