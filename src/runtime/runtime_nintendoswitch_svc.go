package runtime

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
func SvcGetInfo(infoType InfoType, handle uint32, id1 uint64) (value uint64, result Result)

// SvcSetHeapSize  Set the process heap to a given size. It can both extend and shrink the heap.
// size Size of the heap, must be a multiple of 0x2000000 and [2.0.0+] less than 0x18000000.
// Syscall number 0x01.
//go:export svcSetHeapSize
func SvcSetHeapSize(size uint64) (outAddr uintptr, result Result)
