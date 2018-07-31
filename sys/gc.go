package sys

import "runtime/debug"

// Free forces the garbage collector to free memory. Blocks the program until GC is done.
func Free() {
	if MemConservative {
		debug.FreeOSMemory()
	}
}

// ForceFree forces the garbage collector to free memory regardless of sys.MemConservative.
func ForceFree() {
	debug.FreeOSMemory()
}
