// +build nintendoswitch

package runtime

import (
	"unsafe"
)

// Result svcOutputDebugString(const char *str, u64 size)
//go:export svcOutputDebugString
func _SvcOutputDebugString(str unsafe.Pointer, size uint64) uint64

func OutputDebugChar(v byte) {
	_SvcOutputDebugString(unsafe.Pointer(&v), 1)
}

//go:export malloc
func extalloc(size uintptr) unsafe.Pointer

//export free
func extfree(ptr unsafe.Pointer)
