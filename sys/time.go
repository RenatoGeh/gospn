package sys

import (
	"time"
)

var (
	st time.Time
)

// StartTimer manually starts the global timer.
func StartTimer() {
	st = time.Now()
}

// StopTimer manually stops the global timer. Returns the time elapsed since calling StartTimer.
func StopTimer() time.Duration {
	return time.Since(st)
}

// MeasureTime takes a start time s and a string d. MeasureTime will time s with the current time
// and print with GoSPN's Printf a string of format ("%s took %s.", d, s). Consider using it as:
//  func funcToTime() {
//    defer sys.MeasureTime(time.Now(), "funcToTime")
//    // ...
//  }
// Since arguments are evaluated before function call, the function will print the correct
// funcToTime run time.
func MeasureTime(s time.Time, d string) {
	e := time.Since(s)
	Printf("%s took %s.\n", d, e)
}
