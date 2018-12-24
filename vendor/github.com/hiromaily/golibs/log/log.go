package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	c "github.com/hiromaily/golibs/color"
	r "github.com/hiromaily/golibs/runtimes"
	u "github.com/hiromaily/golibs/utils"
)

//Output e.g.
//[GOWEB]20:33:18 [DEBUG]conf environment : local
//[GOWEB]2016/05/17 20:09:32 log.go:132: [DEBUG]conf environment : local

// LogStatus is logStatus
type LogStatus uint8

const (
	// DebugStatus is at debug level
	DebugStatus LogStatus = iota + 1
	// InfoStatus is at info level
	InfoStatus
	// WarningStatus is at warning level
	WarningStatus
	// ErrorStatus is at error level
	ErrorStatus
	// FatalStatus is at fatal level
	FatalStatus
	// StackStatus is at stack level
	StackStatus
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
	// StackPrefix is of fatal prefix
	StackPrefix string = "[STACK]"
)

// LogFmt is log format
type LogFmt int

const (
	NoDateNoFile      LogFmt = 0
	OnlyTime          LogFmt = log.Ltime
	TimeShortFile     LogFmt = log.Ltime | log.Lshortfile //should be default
	DateTimeShortFile LogFmt = log.LstdFlags | log.Lshortfile
)

// Int is to convert to type int
func (lf LogFmt) Int() int {
	return int(lf)
}

/*
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger

*/

// LogType is console or file to output
type LogType uint8

const (
	File LogType = iota + 1
	Console
)

// Logger is for log object
type Logger struct {
	logger    *log.Logger
	logLevel  LogStatus
	logType   LogType
	separator string
}

var (
	logger    *log.Logger
	logLevel  LogStatus = 1
	logType   LogType
	separator string
)

//-----------------------------------------------------------------------------
// functions
//-----------------------------------------------------------------------------

func getStatus(key string) LogStatus {
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
	case "Stack":
		return StackStatus
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
	case "Stack":
		return StackPrefix
	default:
		return ""
	}
}

// setColor is to set color
func setColor(key, val string) string {
	switch key {
	case "Debug", "Debugf":
		return c.Add(c.SkyBlue, val)
	case "Info", "Infof":
		return c.Add(c.Green, val)
	case "Warn", "Warnf":
		return c.Add(c.Yellow, val)
	case "Error", "Errorf":
		return c.Add(c.DeepPink, val)
	case "Fatal", "Fatalf":
		return c.Add(c.Red, val)
	default:
		return val
	}
}

// openFile is for output log file
func openFile(logger *log.Logger, fileName string) {
	if fileName == "" {
		return
	}
	//ファイルがなければ作成する
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening file :", err.Error())
	}
	logger.SetOutput(f)
}

// currentFunc is to get current func name
func currentFunc(skip int) string {
	programCounter, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	sl := strings.Split(runtime.FuncForPC(programCounter).Name(), ".")
	return sl[len(sl)-1]
}

//-----------------------------------------------------------------------------
// logger instance
//-----------------------------------------------------------------------------

// New is to create new log object
func New(level LogStatus, logFmt LogFmt, prefix, fileName, delimiter string) *Logger {
	logConf := &Logger{}
	logConf.logLevel = level
	logConf.separator = delimiter

	logConf.logger = log.New(os.Stderr, prefix, logFmt.Int())

	//Log File Path
	if fileName == "" {
		logConf.logType = Console
	} else {
		logConf.logType = File

		openFile(logConf.logger, fileName)
	}

	return logConf
}

// Out is to output log
func (lo *Logger) Out(key string, v ...interface{}) {
	nv := u.Unshift(v, getPrefix(key))

	if lo.logLevel <= getStatus(key) {
		if lo.logType == File {
			lo.logger.Output(3, fmt.Sprint(nv...))
		} else {
			lo.logger.Output(3, setColor(key, fmt.Sprint(nv...)))
		}
	}
}

// Outf is to output log with format
func (lo *Logger) Outf(key, format string, v ...interface{}) {
	if lo.logLevel <= getStatus(key) {
		if lo.logType == File {
			lo.logger.Output(3, fmt.Sprintf(getPrefix(key)+format, v...))
		} else {
			lo.logger.Output(3, setColor(key, fmt.Sprintf(getPrefix(key)+format, v...)))
		}
	}
}

// Debug is to call output func for debug log
func (lo *Logger) Debug(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Debugf is to call output func for debug log with format
func (lo *Logger) Debugf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

// Info is to call output func for info log
func (lo *Logger) Info(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Infof is to call output func for info log with format
func (lo *Logger) Infof(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

// Warn is to call output func for warn log
func (lo *Logger) Warn(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Warnf is to call output func for warn log with format
func (lo *Logger) Warnf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

// Error is to call output func for error log
func (lo *Logger) Error(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Errorf is to call output func for error log with format
func (lo *Logger) Errorf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

// Fatal is to call output func for fatal log
func (lo *Logger) Fatal(v ...interface{}) {
	key := currentFunc(1)
	lo.Out(key, v...)
}

// Fatalf is to call output func for fatal log with format
func (lo *Logger) Fatalf(format string, v ...interface{}) {
	key := currentFunc(1)
	lo.Outf(key, format, v...)
}

// Stack is to call output func for stack trace
func (lo *Logger) Stack() {
	info := r.GetStackTrace(lo.separator)
	msg := "\n"
	for i := len(info) - 2; i > 0; i-- {
		v := info[i]
		msg += fmt.Sprintf("%02d: [Function]%s [File]%s:%d\n", i, v.FunctionName, v.FileName, v.FileLine)
	}

	key := currentFunc(1)
	lo.Out(key, msg)
}

//-----------------------------------------------------------------------------
// singleton object
//-----------------------------------------------------------------------------

// InitializeLog is to initialize base log object using default setting
func InitializeLog(level LogStatus, logFmt LogFmt, prefix, fileName, delimiter string) {
	logLevel = level
	separator = delimiter

	//Log Object
	logger = log.New(os.Stderr, prefix, logFmt.Int())

	//Log File Path
	if fileName == "" {
		logType = Console
	} else {
		logType = File
		openFile(logger, fileName)
	}
}

// out is to output log
func out(key string, v ...interface{}) {
	nv := u.Unshift(v, getPrefix(key))

	if logLevel <= getStatus(key) {
		if logType == File {
			//file
			logger.Output(3, fmt.Sprint(nv...))
		} else {
			logger.Output(3, setColor(key, fmt.Sprint(nv...)))
		}
	}
}

// outf is to output log with format
func outf(key, format string, v ...interface{}) {
	if logLevel <= getStatus(key) {
		if logType == File {
			//file
			logger.Output(3, fmt.Sprintf(getPrefix(key)+format, v...))
		} else {
			logger.Output(3, setColor(key, fmt.Sprintf(getPrefix(key)+format, v...)))
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

// Stack is to call output func for stack trace
func Stack() {
	info := r.GetStackTrace(separator)
	msg := "\n"
	for i := len(info) - 2; i > 0; i-- {
		v := info[i]
		msg += fmt.Sprintf("%02d: [Function]%s [File]%s:%d\n", i, v.FunctionName, v.FileName, v.FileLine)
	}

	key := currentFunc(1)
	out(key, msg)
}
