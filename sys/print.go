package sys

import (
	"fmt"
)

// Printf is a wrapper for fmt.Printf. Only prints if Verbose is set to true.
func Printf(str string, vals ...interface{}) {
	if !Verbose {
		return
	}
	if len(vals) == 0 {
		fmt.Printf(str)
	} else {
		fmt.Printf(str, vals...)
	}
}

// Println is a wrapper for fmt.Println. Only prints if Verbose is set to true.
func Println(str string) {
	if Verbose {
		fmt.Println(str)
	}
}
