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
	// DebugStatus is at debug level
	DebugStatus uint8 = iota + 1
	// InfoStatus is at info level
	InfoStatus
	// WarningStatus is at warning level
	WarningStatus
	// ErrorStatus is at error level
	ErrorStatus
	// FatalStatus is at fatal level
	FatalStatus
	// LogOff is at no log
	LogOff
)

const (
	// DebugPrefix is of debug prefix
	DebugPrefix string = "[DEBUG]"
	// InfoPrefix is of info prefix
	InfoPrefix string = "[INFO]"
	// WarningPrefix is of warning prefix
	WarningPrefix string = "[WARNING]"
	// ErrorPrefix is of error prefix
	ErrorPrefix string = "[ERROR]"
	// FatalPrefix is of fatal prefix
	FatalPrefix string = "[FATAL]"
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

// Object is for log object
type Object struct {
	loggerStd    *log.Logger
	loggerFile   *log.Logger
	logLevel     uint8
	logFileLevel uint8
}

var (
	logLevel     uint8 = 1
	logFileLevel uint8 = 4
	filePathName       = "/var/log/go/xxxx.log"
	logStdOut    *log.Logger
	logFileOut   *log.Logger
)

// currentFunc is to get current func name
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
		return DebugStatus
	case "Info", "Infof":
		return InfoStatus
	case "Warn", "Warnf":
		return WarningStatus
	case "Error", "Errorf":
		return ErrorStatus
	case "Fatal", "Fatalf":
		return FatalStatus
	default:
		return 0
	}
}

func getPrefix(key string) string {
	switch key {
	case "Debug", "Debugf":
		return DebugPrefix
	case "Info", "Infof":
		return InfoPrefix
	case "Warn", "Warnf":
		return WarningPrefix
	case "Error", "Errorf":
		return ErrorPrefix
	case "Fatal", "Fatalf":
		return FatalPrefix
	default:
		return ""
	}
}

// openFile is for output log file
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

// setColor is to set color
// TODO: work in progress
func setColor() {
	logStdOut.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
}

// New is to create new log object
func New(level, fileLevel uint8, logFmt int, prefix, fileName string) *Object {
	logObj := &Object{}
	logObj.logLevel = level
	logObj.logFileLevel = fileLevel

	//Log File Path
	if fileName == "" {
		fileName = filePathName
	}

	//Log Format
	if logFmt == 0 {
		//date and file is not shown.
	} else if logFmt == 99 {
		//logFmt = log.Ltime                               //2
		logFmt = log.Ltime | log.Lshortfile //18
		//logFmt = log.LstdFlags | log.Lshortfile          //19
		//logFmt = log.Ldate | log.Ltime | log.Lshortfile  //19
	}

	logObj.loggerStd = log.New(os.Stderr, prefix, logFmt)
	logObj.loggerFile = log.New(os.Stderr, prefix, logFmt)

	if fileLevel != LogOff {
		openFile(logObj.loggerFile, fileName)
	}

	return logObj
}

// Out is to output log
func (lo *Object) Out(key string, v ...interface{}) {
	nv := u.Unshift(v, getPrefix(key))

	if lo.logLevel <= getStatus(key) {
		if lo.logFileLevel <= getStatus(key) {
			lo.loggerFile.Output(2, fmt.Sprint(nv...))
		} else {
			lo.loggerStd.Output(2, fmt.Sprint(nv...))
		}
	}
}

// Outf is to output log with format
func (lo *Object) Outf(key, format string, v ...interface{}) {
	if lo.logLevel <= getStatus(key) {
		if lo.logFileLevel <= getStatus(key) {
			lo.loggerFile.Output(2, fmt.Sprintf(getPrefix(key)+format, v...))
		} else {
			lo.loggerStd.Output(2, fmt.Sprintf(getPrefix(key)+format, v...))
		}
	}
}

// Debug is to call output func for debug log
func (lo *Object) Debug(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Debugf is to call output func for debug log with format
func (lo *Object) Debugf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

// Info is to call output func for info log
func (lo *Object) Info(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Infof is to call output func for info log with format
func (lo *Object) Infof(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

// Warn is to call output func for warn log
func (lo *Object) Warn(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Warnf is to call output func for warn log with format
func (lo *Object) Warnf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

// Error is to call output func for error log
func (lo *Object) Error(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Errorf is to call output func for error log with format
func (lo *Object) Errorf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

// Fatal is to call output func for fatal log
func (lo *Object) Fatal(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Fatalf is to call output func for fatal log with format
func (lo *Object) Fatalf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

//-----------------------------------------------------------------------------
// singleton object
//-----------------------------------------------------------------------------

// InitializeLog is to initialize base log object using default setting
func InitializeLog(level, fileLevel uint8, logFmt int, prefix, fileName string) {
	logLevel = level
	logFileLevel = fileLevel

	//Log File Path
	if fileName == "" {
		fileName = filePathName
	}

	//Log Format
	if logFmt == 0 {
		//date and file is not shown.
	} else if logFmt == 99 {
		//logFmt = log.Ltime                               //2
		logFmt = log.Ltime | log.Lshortfile //18 (best for me)
		//logFmt = log.LstdFlags | log.Lshortfile          //19
		//logFmt = log.Ldate | log.Ltime | log.Lshortfile  //19
	}

	//Log Object
	logStdOut = log.New(os.Stderr, prefix, logFmt)
	// color mode
	setColor()

	logFileOut = log.New(os.Stderr, prefix, logFmt)
	if fileLevel != LogOff {
		openFile(logFileOut, fileName)
	}
}

// out is to output log
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

// outf is to output log with format
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

// Debug is to call output func for debug log
func Debug(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

// Debugf is to call output func for debug log with format
func Debugf(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}

// Info is to call output func for info log
func Info(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

// Infof is to call output func for info log with format
func Infof(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}

// Warn is to call output func for warn log
func Warn(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

// Warnf is to call output func for warn log with format
func Warnf(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}

// Error is to call output func for error log
func Error(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

// Errorf is to call output func for error log with format
func Errorf(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}

// Fatal is to call output func for fatal log
func Fatal(v ...interface{}) {
	key := currentFunc(1)
	out(key, v...)
}

// Fatalf is to call output func for fatal log with format
func Fatalf(format string, v ...interface{}) {
	key := currentFunc(1)
	outf(key, format, v...)
}
