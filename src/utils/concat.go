package utils

import (
	"bytes"
)

func StringConcat(s1, s2 string) string {
	var buffer bytes.Buffer
	buffer.WriteString(s1)
	buffer.WriteString(s2)
	return buffer.String()
}
