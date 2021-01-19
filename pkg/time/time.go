package time

import (
	"fmt"
	"time"
)

// Track is to track elapsed time
//  e.g. Caller: defer Track(time.Now(), "parseFile()")
//  https://medium.com/@2xb/execution-time-tracking-in-golang-9379aebfe20e#.ffxgxejim
func Track(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
