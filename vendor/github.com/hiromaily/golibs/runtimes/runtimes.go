package runtimes

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	re = regexp.MustCompile(`^(\S.+)\.(\S.+)$`)
)

type CallerInfo struct {
	PackageName  string
	FunctionName string
	FileName     string
	FileLine     int
}

func formatPath(path, separator string) string {
	ret := strings.Split(path, separator)
	if len(ret) > 1 {
		return fmt.Sprintf(".%s", ret[len(ret)-1])
	}
	return path
}

func dumpStackTrace(separator string) (callerInfo []*CallerInfo) {
	for i := 1; ; i++ {
		pc, _, _, ok := runtime.Caller(i) // https://golang.org/pkg/runtime/#Caller
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		fileName, fileLine := fn.FileLine(pc)

		// format path
		if separator != "" {
			fileName = formatPath(fileName, separator)
		}

		additionalInfo := re.FindStringSubmatch(fn.Name())
		callerInfo = append(callerInfo, &CallerInfo{
			PackageName:  additionalInfo[1],
			FunctionName: additionalInfo[2],
			FileName:     fileName,
			FileLine:     fileLine,
		})
	}
	return callerInfo[1:]
}

func GetOS() string {
	return runtime.GOOS
}

func GetStackTrace(separator string) []*CallerInfo {
	info := dumpStackTrace(separator)
	return info
}

func TraceAllHistory(w io.Writer, separator string) {
	info := dumpStackTrace(separator)
	for i := len(info) - 1; i > -1; i-- {
		v := info[i]
		//fmt.Printf("%02d: %s%s@%s:%d\n", i, v.PackageName, v.FunctionName, v.FileName, v.FileLine)
		fmt.Fprintf(w, "%02d: [Function]%s [File]%s:%d\n", i, v.FunctionName, v.FileName, v.FileLine)
	}
}

// CurrentFunc is to get current func name
func CurrentFunc(skip int) string {
	programCounter, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	sl := strings.Split(runtime.FuncForPC(programCounter).Name(), ".")
	return sl[len(sl)-1]
}

// CurrentFuncV2 is to get current func name
func CurrentFuncV2() []byte {
	b := make([]byte, 250)
	b = b[:runtime.Stack(b, false)]
	for i := 0; i < 3; i++ {
		j := bytes.IndexByte(b, '\n')
		if j < 0 {
			return nil
		}

		b = b[j+1:]
	}
	i := bytes.IndexByte(b, '(')
	if i < 0 {
		return nil
	}

	return b[:i]
}

func DebugStack() {
	// Stack returns a formatted stack trace of the goroutine that calls it.
	// It calls runtime.Stack with a large enough buffer to capture the entire trace.
	fmt.Println(string(debug.Stack()))
}

func DebugPrintStack() {
	// PrintStack prints to standard error the stack trace returned by runtime.Stack.
	debug.PrintStack()
}
