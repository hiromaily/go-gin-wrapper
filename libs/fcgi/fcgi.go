package fcgi

import (
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http/fcgi"
	"os"
)

//func IsDebugging() bool {
//	return ginMode == debugCode
//}

func debugPrint(format string, values ...interface{}) {
	//if IsDebugging() {
	//	log.Printf("[GIN-debug] "+format, values...)
	//}
	log.Printf("[GIN-debug] "+format, values...)
}

func debugPrintError(err error) {
	if err != nil {
		debugPrint("[ERROR] %v\n", err)
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); len(port) > 0 {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func Run(engine *gin.Engine, addr ...string) (err error) {
	defer func() { debugPrintError(err) }()

	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)
	//err = http.ListenAndServe(address, engine)
	ltn, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	err = fcgi.Serve(ltn, engine)

	return
}
