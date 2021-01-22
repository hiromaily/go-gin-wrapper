package time

import (
	"fmt"
	"time"
)

// Track tracks elapsed time
//  e.g. Caller: defer Track(time.Now(), "parseFile()")
func Track(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
