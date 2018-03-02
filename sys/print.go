package sys

import (
	"fmt"
	"os"
)

var (
	fout = os.Stdout
)

// SetLogFile forces GoSPN to write all content to a log file. If filename is empty, writes to
// standard output. Returns the file pointer and errors if any.
func SetLogFile(filename string) (*os.File, error) {
	if filename == "" {
		fout = os.Stdout
		return nil, nil
	}
	f, err := os.Create(filename)
	return f, err
}

// LogFile returns the log file. The default log file is os.Stdout, which is the standard output
// device.
func LogFile() *os.File {
	return fout
}

// Printf is a wrapper for fmt.Printf. Only prints if Verbose is set to true.
func Printf(str string, vals ...interface{}) {
	if !Verbose {
		return
	}
	if len(vals) == 0 {
		fmt.Fprintf(fout, str)
	} else {
		fmt.Fprintf(fout, str, vals...)
	}
}

// Println is a wrapper for fmt.Println. Only prints if Verbose is set to true.
func Println(str string) {
	if Verbose {
		fmt.Fprintln(fout, str)
	}
}
