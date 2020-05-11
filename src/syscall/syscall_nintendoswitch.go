// +build nintendoswitch

package syscall

import "unsafe"

// InfoType Types of information to use on SvcGetInfo System Call
type InfoType uint32

// Result represents a result state from a System Call
type Result uint32

// Permission represents memory permission bitmasks.
type Permission uint32

//
// Nintendo Switch - Horizon System Calls
// https://switchbrew.org/wiki/SVC
//

const (
	// Pseudo handle representing an invalid object
	INVALID_HANDLE = 0

	// Pseudo handle for the current process.
	CUR_PROCESS_HANDLE = 0xFFFF8001

	// Pseudo handle for the current thread.
	CUR_THREAD_HANDLE = 0xFFFF8000
)

/// Memory information structure
type MemoryInfo struct {
	Addr           uint64 //< Base address.
	Size           uint64 //< Size.
	MemType        uint32 //< Memory type (see lower 8 bits of \ref MemoryState).
	Attr           uint32 //< Memory attributes (see \ref MemoryAttribute).
	Perm           uint32 //< Memory permissions (see \ref Permission).
	DeviceRefcount uint32 //< Device reference count.
	IpcRefcount    uint32 //< IPC reference count.
	Padding        uint32 //< Padding.
}

type PhysicalMemoryInfo struct {
	PhysicalAddress uint64 //< Physical address.
	VirtualAddress  uint64 //< Virtual address.
	Size            uint64 //< Size.
}

// SvcGetInfo Retrieves information about the system, or a certain kernel object.
// infoType First ID of the property to retrieve
// handle Handle of the object to retrieve information from, or INVALID_HANDLE to retrieve information about the system.
// Syscall number 0x29
//go:export svcGetInfo
func _svcGetInfo(outPtr unsafe.Pointer, infoType InfoType, handle uint32, id1 uint64) Result

// SvcSetHeapSize  Set the process heap to a given size. It can both extend and shrink the heap.
// size Size of the heap, must be a multiple of 0x2000000 and [2.0.0+] less than 0x18000000.
// Syscall number 0x01.
//go:export svcSetHeapSize
func SvcSetHeapSize(size uint64) (outAddr uintptr, result Result)

// SvcSetMemoryPermission Set the memory permissions of a (page-aligned) range of memory.
// Perm_X is not allowed. Setting write-only is not allowed either (Perm_W).
// This can be used to move back and forth between Perm_None, Perm_R and Perm_Rw.
// Syscall number 0x01
//go:export svcSetMemoryPermission
func SvcSetMemoryPermission(addr unsafe.Pointer, size uint64, perm Permission) Result
