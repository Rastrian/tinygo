// +build nintendoswitch

package runtime

import (
	"reflect"
	"unsafe"
)

//go:nobounds
func printstring(s string) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	_SvcOutputDebugString(unsafe.Pointer(sh.Data), uint64(sh.Len))
}
