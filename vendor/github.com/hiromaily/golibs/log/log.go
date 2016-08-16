package log

import (
	"fmt"
	u "github.com/hiromaily/golibs/utils"
	"github.com/shiena/ansicolor"
	"log"
	"os"
	"runtime"
	"strings"
)

//Output e.g.
//[GOWEB]20:33:18 [DEBUG]conf environment : local
//[GOWEB]2016/05/17 20:09:32 log.go:132: [DEBUG]conf environment : local

const (
	DEBUG_STATUS uint8 = iota + 1
	INFO_STATUS
	WARNING_STATUS
	ERROR_STATUS
	FATAL_STATUS
	LOG_OFF_COUNT
)

const (
	DEBUG_PREFIX   string = "[DEBUG]"
	INFO_PREFIX    string = "[INFO]"
	WARNING_PREFIX string = "[WARNING]"
	ERROR_PREFIX   string = "[ERROR]"
	FATAL_PREFIX   string = "[FATAL]"
)

/*
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger

*/

var (
	logLevel     uint8  = 1
	logFileLevel uint8  = 4
	filePathName string = "/var/log/go/xxxx.log"
)

type LogObject struct {
	loggerStd    *log.Logger
	loggerFile   *log.Logger
	logLevel     uint8
	logFileLevel uint8
}

var (
	logStdOut  *log.Logger
	logFileOut *log.Logger
)

// get current func name
func currentFunc(skip int) string {
	programCounter, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	sl := strings.Split(runtime.FuncForPC(programCounter).Name(), ".")
	return sl[len(sl)-1]
}

func getStatus(key string) uint8 {
	switch key {
	case "Debug", "Debugf":
		return DEBUG_STATUS
	case "Info", "Infof":
		return INFO_STATUS
	case "Warn", "Warnf":
		return WARNING_STATUS
	case "Error", "Errorf":
		return ERROR_STATUS
	case "Fatal", "Fatalf":
		return FATAL_STATUS
	default:
		return 0
	}
}

func getPrefix(key string) string {
	switch key {
	case "Debug", "Debugf":
		return DEBUG_PREFIX
	case "Info", "Infof":
		return INFO_PREFIX
	case "Warn", "Warnf":
		return WARNING_PREFIX
	case "Error", "Errorf":
		return ERROR_PREFIX
	case "Fatal", "Fatalf":
		return FATAL_PREFIX
	default:
		return ""
	}
}

//for output log file
func openFile(logger *log.Logger, fileName string) {
	if fileName == "" {
		return
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening file :", err.Error())
	}
	logger.SetOutput(f)
}

// Set color
func setColor() {
	logStdOut.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
}

//Create New Original Object
func New(level, fileLevel uint8, logFmt int, prefix, fileName string) *LogObject {
	logObj := &LogObject{}
	logObj.logLevel = level
	logObj.logFileLevel = fileLevel

	//Log File Path
	if fileName == "" {
		fileName = filePathName
	}

	//Log Format
	if logFmt == 0 {
		//logFmt = log.Ltime                               //2
		logFmt = log.Ltime | log.Lshortfile //18
		//logFmt = log.LstdFlags | log.Lshortfile          //19
		//logFmt = log.Ldate | log.Ltime | log.Lshortfile  //19
	}

	logObj.loggerStd = log.New(os.Stderr, prefix, logFmt)
	logObj.loggerFile = log.New(os.Stderr, prefix, logFmt)

	if fileLevel != LOG_OFF_COUNT {
		openFile(logObj.loggerFile, fileName)
	}

	return logObj
}

func (lo *LogObject) Out(key string, v ...interface{}) {
	nv := u.Unshift(v, getPrefix(key))

	if lo.logLevel <= getStatus(key) {
		if lo.logFileLevel <= getStatus(key) {
			lo.loggerFile.Output(2, fmt.Sprint(nv...))
		} else {
			lo.loggerStd.Output(2, fmt.Sprint(nv...))
		}
	}
}

func (lo *LogObject) Outf(key, format string, v ...interface{}) {
	if lo.logLevel <= getStatus(key) {
		if lo.logFileLevel <= getStatus(key) {
			lo.loggerFile.Output(2, fmt.Sprintf(getPrefix(key)+format, v...))
		} else {
			lo.loggerStd.Output(2, fmt.Sprintf(getPrefix(key)+format, v...))
		}
	}
}

func (lo *LogObject) Debug(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

func (lo *LogObject) Debugf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

func (lo *LogObject) Info(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

func (lo *LogObject) Infof(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

func (lo *LogObject) Warn(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

func (lo *LogObject) Warnf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

func (lo *LogObject) Error(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

func (lo *LogObject) Errorf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

func (lo *LogObject) Fatal(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

func (lo *LogObject) Fatalf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

//Initialize base log object using default setting
func InitializeLog(level, fileLevel uint8, logFmt int, prefix, fileName string) {
	logLevel = level
	logFileLevel = fileLevel

	//Log File Path
	if fileName == "" {
		fileName = filePathName
	}

	//Log Format
	if logFmt == 0 {
		//logFmt = log.Ltime                               //2
		logFmt = log.Ltime | log.Lshortfile //18
		//logFmt = log.LstdFlags | log.Lshortfile          //19
		//logFmt = log.Ldate | log.Ltime | log.Lshortfile  //19
	}

	//Log Object
	logStdOut = log.New(os.Stderr, prefix, logFmt)
	// color mode
	setColor()

	logFileOut = log.New(os.Stderr, prefix, logFmt)
	if fileLevel != LOG_OFF_COUNT {
		openFile(logFileOut, fileName)
	}
}

func out(key string, v ...interface{}) {
	nv := u.Unshift(v, getPrefix(key))

	if logLevel <= getStatus(key) {
		if logFileLevel <= getStatus(key) {
			//file
			logFileOut.Output(3, fmt.Sprint(nv...))
		} else {
			logStdOut.Output(3, fmt.Sprint(nv...))
		}
	}
}

func outf(key, format string, v ...interface{}) {
	if logLevel <= getStatus(key) {
		if logFileLevel <= getStatus(key) {
			//file
			logFileOut.Output(3, fmt.Sprintf(getPrefix(key)+format, v...))
		} else {
			logStdOut.Output(3, fmt.Sprintf(getPrefix(key)+format, v...))
		}
	}
}

//Debug
func Debug(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

func Debugf(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}

//Info
func Info(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

func Infof(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}

//Warn
func Warn(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

func Warnf(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}

//Error
func Error(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

func Errorf(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}

//Fatal
func Fatal(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

func Fatalf(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}
