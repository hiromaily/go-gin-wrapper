package fcgi

import (
	"log"
	"net"
	"net/http/fcgi"

	"github.com/gin-gonic/gin"
)

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func Run(engine *gin.Engine, addr string) error {
	log.Printf("[GIN-debug] Listening and serving HTTP on %s\n", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return fcgi.Serve(listener, engine)
}
