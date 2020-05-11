// +build nintendoswitch

package runtime

import (
	"unsafe"
)

//
// Nintendo Switch - Horizon System Calls
// https://switchbrew.org/wiki/SVC
//

const (
	// Pseudo handle representing an invalid object
	InvalidHandle = 0

	// Pseudo handle for the current process.
	CurProcessHandle = 0xFFFF8001

	// Pseudo handle for the current thread.
	CurThreadHandle = 0xFFFF8000

	// Maximum number of objects that can be waited on by \ref svcWaitSynchronization (Horizon kernel limitation).
	MaxWaitObjects = 0x40
)

// InfoType Types of information to use on SvcGetInfo System Call
type InfoType uint32

// Result represents a result state from a System Call
type Result uint32

// Predefined Results
const (
	ResultOK Result = 0
)

// GetInfo IDs.
const (
	InfoType_CoreMask                    InfoType = 0          // Bitmask of allowed Core IDs.
	InfoType_PriorityMask                InfoType = 1          // Bitmask of allowed Thread Priorities.
	InfoType_AliasRegionAddress          InfoType = 2          // Base of the Alias memory region.
	InfoType_AliasRegionSize             InfoType = 3          // Size of the Alias memory region.
	InfoType_HeapRegionAddress           InfoType = 4          // Base of the Heap memory region.
	InfoType_HeapRegionSize              InfoType = 5          // Size of the Heap memory region.
	InfoType_TotalMemorySize             InfoType = 6          // Total amount of memory available for process.
	InfoType_UsedMemorySize              InfoType = 7          // Amount of memory currently used by process.
	InfoType_DebuggerAttached            InfoType = 8          // Whether current process is being debugged.
	InfoType_ResourceLimit               InfoType = 9          // Current process's resource limit handle.
	InfoType_IdleTickCount               InfoType = 10         // Number of idle ticks on CPU.
	InfoType_RandomEntropy               InfoType = 11         // [2.0.0+] Random entropy for current process.
	InfoType_AslrRegionAddress           InfoType = 12         // [2.0.0+] Base of the process's address space.
	InfoType_AslrRegionSize              InfoType = 13         // [2.0.0+] Size of the process's address space.
	InfoType_StackRegionAddress          InfoType = 14         // [2.0.0+] Base of the Stack memory region.
	InfoType_StackRegionSize             InfoType = 15         // [2.0.0+] Size of the Stack memory region.
	InfoType_SystemResourceSizeTotal     InfoType = 16         // [3.0.0+] Total memory allocated for process memory management.
	InfoType_SystemResourceSizeUsed      InfoType = 17         // [3.0.0+] Amount of memory currently used by process memory management.
	InfoType_ProgramId                   InfoType = 18         // [3.0.0+] Program ID for the process.
	InfoType_InitialProcessIdRange       InfoType = 19         // [4.0.0-4.1.0] Min/max initial process IDs.
	InfoType_UserExceptionContextAddress InfoType = 20         // [5.0.0+] Address of the process's exception context (for break).
	InfoType_TotalNonSystemMemorySize    InfoType = 21         // [6.0.0+] Total amount of memory available for process, excluding that for process memory management.
	InfoType_UsedNonSystemMemorySize     InfoType = 22         // [6.0.0+] Amount of memory used by process, excluding that for process memory management.
	InfoType_IsApplication               InfoType = 23         // [9.0.0+] Whether the specified process is an Application.
	InfoType_ThreadTickCount             InfoType = 0xF0000002 // Number of ticks spent on thread.
)

// SvcGetInfo Retrieves information about the system, or a certain kernel object.
// infoType First ID of the property to retrieve
// handle Handle of the object to retrieve information from, or INVALID_HANDLE to retrieve information about the system.
// Syscall number 0x29
//go:export svcGetInfo
func _SvcGetInfo(outPtr unsafe.Pointer, infoType InfoType, handle uint32, id1 uint64) Result

// SvcSetHeapSize  Set the process heap to a given size. It can both extend and shrink the heap.
// size Size of the heap, must be a multiple of 0x2000000 and [2.0.0+] less than 0x18000000.
// Syscall number 0x01.
//go:export svcSetHeapSize
func SvcSetHeapSize(size uint64) (outAddr uintptr, result Result)

func SvcGetInfo(infoType InfoType, handle uint32, id1 uint64) (value uint64, result Result) {
	r := _SvcGetInfo(unsafe.Pointer(&value), infoType, handle, id1)
	return value, r
}

// Result svcOutputDebugString(const char *str, u64 size)
//go:export svcOutputDebugString
func _SvcOutputDebugString(str unsafe.Pointer, size uint64) Result

func OutputDebugChar(v byte) {
	_SvcOutputDebugString(unsafe.Pointer(&v), 1)
}

// Result svcBreak(u32 breakReason, u64 inval1, u64 inval2);
//go:export svcBreak
func svcBreak(reason uint32, a, b uint64)

//go:export consoleInit
func consoleInit(c unsafe.Pointer)

//go:export printf
func libc_printf(char unsafe.Pointer, val uint64)

// (int fd, const void *buf, size_t cnt)
//go:export write
func libc_write(fd int, buffer unsafe.Pointer, size uint64) int

func printf(data string, val uint64) {
	libc_printf(unsafe.Pointer((*_string)(unsafe.Pointer(&data)).ptr), val)
}

//go:export appletMainLoop
func appletMainLoop() bool

//go:export consoleUpdate
func consoleUpdate(c unsafe.Pointer)

//go:export consoleExit
func consoleExit(c unsafe.Pointer)

//go:export hidScanInput
func hidScanInput()

//u64 hidKeysDown(HidControllerID id);
//go:export hidKeysDown
func hidKeysDown(id uint64) uint64

// void *malloc(size_t size);
//go:export malloc
func libc_malloc(size uint64) unsafe.Pointer

//go:export __nx_exit
func __nx_exit(code int)
