package signal

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

// GOTRACEBACK=single
// GOTRACEBACK=all

// StartSignal is to wait signal using goroutine
func StartSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int)
	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGHUP xxx
			case syscall.SIGHUP:
				fmt.Println("[Signal] hungup")

			// kill -SIGINT xxx or Ctrl+c
			case syscall.SIGINT:
				fmt.Println("[Signal] interrupt by control + C")
				debug.PrintStack()
				exitChan <- 0

			// kill -SIGTERM xxx
			case syscall.SIGTERM:
				fmt.Println("[Signal] force stop by kill command")
				exitChan <- 0

			// kill -SIGQUIT xxx
			case syscall.SIGQUIT:
				fmt.Println("[Signal] stop and core dump")
				exitChan <- 0

			default:
				fmt.Println("[Signal] Unknown signal.")
				exitChan <- 1
			}
		}
	}()

	<-exitChan
	panic("detected signal") // for all stack trace
}
