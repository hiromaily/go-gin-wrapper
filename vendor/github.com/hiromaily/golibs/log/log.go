package log

import (
	"fmt"
	u "github.com/hiromaily/golibs/utils"
	"log"
	"os"
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
	logger *log.Logger
}

var (
	logStdOut  LogObject = LogObject{}
	logFileOut LogObject = LogObject{}
)

//for output log file
func (self *LogObject) openFile(fileName string) {
	if fileName == "" {
		return
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening file :", err.Error())
	}
	self.logger.SetOutput(f)
}

//Create New Original Object
//e.g.
// lg.New("[ProjectName] ", Ltime|Lshortfile, "/var/log/go/xxx.log")
func New(prefix string, logFmt int, fileName string) (*log.Logger, error) {
	logObj := log.New(os.Stderr, prefix, logFmt)

	if fileName != "" {
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Error opening logfile :", err.Error())
			return nil, err
		}
		//logObj.SetOutput undefined (type *log.Logger has no field or method SetOutput) on golang version 1.4
		logObj.SetOutput(f)
	}

	return logObj, nil
}

//Initialize base log object using default setting
//filePath have to include file name.
//e.g.
// lg.InitializeLog(lg.DEBUG_STATUS, lg.LOG_OFF_COUNT, 0, "[GOWEB]", "/var/log/go/goweb.log")
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
	logStdOut.logger = log.New(os.Stderr, prefix, logFmt)

	logFileOut.logger = log.New(os.Stderr, prefix, logFmt)
	if fileLevel != LOG_OFF_COUNT {
		logFileOut.openFile(fileName)
	}
}

//Debug
func Debug(v ...interface{}) {
	//nv := append([]interface{}{DEBUG_PREFIX}, v...)
	nv := u.Unshift(v, DEBUG_PREFIX)

	if logLevel <= DEBUG_STATUS {
		if logFileLevel <= DEBUG_STATUS {
			//file
			//logFileOut.logger.Print(nv...)
			logFileOut.logger.Output(2, fmt.Sprint(nv...))
		} else {
			//logStdOut.logger.Print(nv...)
			logStdOut.logger.Output(2, fmt.Sprint(nv...))
		}
	}
}

func Debugf(format string, v ...interface{}) {
	if logLevel <= DEBUG_STATUS {
		if logFileLevel <= DEBUG_STATUS {
			//file
			//logFileOut.logger.Printf(DEBUG_PREFIX + format, v...)
			logFileOut.logger.Output(2, fmt.Sprintf(DEBUG_PREFIX+format, v...))
		} else {
			//logStdOut.logger.Printf(DEBUG_PREFIX + format, v...)
			logStdOut.logger.Output(2, fmt.Sprintf(DEBUG_PREFIX+format, v...))
		}
	}
}

//Info
func Info(v ...interface{}) {
	nv := u.Unshift(v, INFO_PREFIX)
	if logLevel <= INFO_STATUS {
		if logFileLevel <= INFO_STATUS {
			//file
			//logFileOut.logger.Print(nv...)
			logFileOut.logger.Output(2, fmt.Sprint(nv...))
		} else {
			//logStdOut.logger.Print(nv...)
			logStdOut.logger.Output(2, fmt.Sprint(nv...))
		}
	}
}

func Infof(format string, v ...interface{}) {
	if logLevel <= INFO_STATUS {
		if logFileLevel <= INFO_STATUS {
			//file
			//logFileOut.logger.Printf(INFO_PREFIX + format, v...)
			logFileOut.logger.Output(2, fmt.Sprintf(INFO_PREFIX+format, v...))
		} else {
			//logStdOut.logger.Printf(INFO_PREFIX + format, v...)
			logStdOut.logger.Output(2, fmt.Sprintf(INFO_PREFIX+format, v...))
		}
	}
}

//Warn
func Warn(v ...interface{}) {
	nv := u.Unshift(v, WARNING_PREFIX)
	if logLevel <= WARNING_STATUS {
		if logFileLevel <= WARNING_STATUS {
			//file
			//logFileOut.logger.Print(nv...)
			logFileOut.logger.Output(2, fmt.Sprint(nv...))
		} else {
			//logStdOut.logger.Print(nv...)
			logStdOut.logger.Output(2, fmt.Sprint(nv...))
		}
	}
}

func Warnf(format string, v ...interface{}) {
	if logLevel <= WARNING_STATUS {
		if logFileLevel <= WARNING_STATUS {
			//file
			//logFileOut.logger.Printf(WARNING_PREFIX + format, v...)
			logFileOut.logger.Output(2, fmt.Sprintf(WARNING_PREFIX+format, v...))
		} else {
			//logStdOut.logger.Printf(WARNING_PREFIX + format, v...)
			logStdOut.logger.Output(2, fmt.Sprintf(WARNING_PREFIX+format, v...))
		}
	}
}

//Error
func Error(v ...interface{}) {
	nv := u.Unshift(v, ERROR_PREFIX)
	if logLevel <= ERROR_STATUS {
		if logFileLevel <= ERROR_STATUS {
			//file
			//logFileOut.logger.Print(nv...)
			logFileOut.logger.Output(2, fmt.Sprint(nv...))
		} else {
			//logStdOut.logger.Print(nv...)
			logStdOut.logger.Output(2, fmt.Sprint(nv...))
		}
	}
}

func Errorf(format string, v ...interface{}) {
	if logLevel <= ERROR_STATUS {
		if logFileLevel <= ERROR_STATUS {
			//file
			//logFileOut.logger.Printf(ERROR_PREFIX + format, v...)
			logFileOut.logger.Output(2, fmt.Sprintf(ERROR_PREFIX+format, v...))
		} else {
			//logStdOut.logger.Printf(ERROR_PREFIX + format, v...)
			logStdOut.logger.Output(2, fmt.Sprintf(ERROR_PREFIX+format, v...))
		}
	}
}

//Fatal
func Fatal(v ...interface{}) {
	nv := u.Unshift(v, FATAL_PREFIX)
	if logLevel <= FATAL_STATUS {
		if logFileLevel <= FATAL_STATUS {
			//file
			//logFileOut.logger.Print(nv...)
			logFileOut.logger.Output(2, fmt.Sprint(nv...))
		} else {
			//logStdOut.logger.Print(nv...)
			logStdOut.logger.Output(2, fmt.Sprint(nv...))
		}
	}
}

func Fatalf(format string, v ...interface{}) {
	if logLevel <= FATAL_STATUS {
		if logFileLevel <= FATAL_STATUS {
			//file
			//logFileOut.logger.Printf(FATAL_PREFIX + format, v...)
			logFileOut.logger.Output(2, fmt.Sprintf(FATAL_PREFIX+format, v...))
		} else {
			//logStdOut.logger.Printf(FATAL_PREFIX + format, v...)
			logStdOut.logger.Output(2, fmt.Sprintf(FATAL_PREFIX+format, v...))
		}
	}
}
