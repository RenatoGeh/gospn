package spn

import (
	"testing"
)

func TestPrint(t *testing.T) {
	S := sampleSPN()
	PrintSPN(S, "/tmp/print_test.spn")
}
