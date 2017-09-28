package sys

import "runtime/debug"

// Free forces the garbage collector to free memory. Blocks the program until GC is done.
func Free() {
	debug.FreeOSMemory()
}
