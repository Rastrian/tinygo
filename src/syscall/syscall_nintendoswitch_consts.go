package syscall

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

// Memory type enumeration (lower 8 bits of \ref MemoryState)
const (
	MemType_Unmapped            = 0x00 // < Unmapped memory.
	MemType_Io                  = 0x01 // < Mapped by kernel capability parsing in \ref svcCreateProcess.
	MemType_Normal              = 0x02 // < Mapped by kernel capability parsing in \ref svcCreateProcess.
	MemType_CodeStatic          = 0x03 // < Mapped during \ref svcCreateProcess.
	MemType_CodeMutable         = 0x04 // < Transition from MemType_CodeStatic performed by \ref svcSetProcessMemoryPermission.
	MemType_Heap                = 0x05 // < Mapped using \ref svcSetHeapSize.
	MemType_SharedMem           = 0x06 // < Mapped using \ref svcMapSharedMemory.
	MemType_WeirdMappedMem      = 0x07 // < Mapped using \ref svcMapMemory.
	MemType_ModuleCodeStatic    = 0x08 // < Mapped using \ref svcMapProcessCodeMemory.
	MemType_ModuleCodeMutable   = 0x09 // < Transition from \ref MemType_ModuleCodeStatic performed by \ref svcSetProcessMemoryPermission.
	MemType_IpcBuffer0          = 0x0A // < IPC buffers with descriptor flags=0.
	MemType_MappedMemory        = 0x0B // < Mapped using \ref svcMapMemory.
	MemType_ThreadLocal         = 0x0C // < Mapped during \ref svcCreateThread.
	MemType_TransferMemIsolated = 0x0D // < Mapped using \ref svcMapTransferMemory when the owning process has perm=0.
	MemType_TransferMem         = 0x0E // < Mapped using \ref svcMapTransferMemory when the owning process has perm!=0.
	MemType_ProcessMem          = 0x0F // < Mapped using \ref svcMapProcessMemory.
	MemType_Reserved            = 0x10 // < Reserved.
	MemType_IpcBuffer1          = 0x11 // < IPC buffers with descriptor flags=1.
	MemType_IpcBuffer3          = 0x12 // < IPC buffers with descriptor flags=3.
	MemType_KernelStack         = 0x13 // < Mapped in kernel during \ref svcCreateThread.
	MemType_CodeReadOnly        = 0x14 // < Mapped in kernel during \ref svcControlCodeMemory.
	MemType_CodeWritable        = 0x15 // < Mapped in kernel during \ref svcControlCodeMemory.
)

/// Memory state bitmasks.
const (
	MemState_Type                       = 0xFF                     // < Type field (see \ref MemoryType).
	MemState_PermChangeAllowed          = 1 << 8                   // < Permission change allowed.
	MemState_ForceRwByDebugSyscalls     = 1 << 9                   // < Force read/writable by debug syscalls.
	MemState_IpcSendAllowed_Type0       = 1 << 10                  // < IPC type 0 send allowed.
	MemState_IpcSendAllowed_Type3       = 1 << 11                  // < IPC type 3 send allowed.
	MemState_IpcSendAllowed_Type1       = 1 << 12                  // < IPC type 1 send allowed.
	MemState_ProcessPermChangeAllowed   = 1 << 14                  // < Process permission change allowed.
	MemState_MapAllowed                 = 1 << 15                  // < Map allowed.
	MemState_UnmapProcessCodeMemAllowed = 1 << 16                  // < Unmap process code memory allowed.
	MemState_TransferMemAllowed         = 1 << 17                  // < Transfer memory allowed.
	MemState_QueryPAddrAllowed          = 1 << 18                  // < Query physical address allowed.
	MemState_MapDeviceAllowed           = 1 << 19                  // < Map device allowed (\ref svcMapDeviceAddressSpace and \ref svcMapDeviceAddressSpaceByForce.
	MemState_MapDeviceAlignedAllowed    = 1 << 20                  // < Map device aligned allowed.
	MemState_IpcBufferAllowed           = 1 << 21                  // < IPC buffer allowed.
	MemState_IsPoolAllocated            = 1 << 22                  // < Is pool allocated.
	MemState_IsRefCounted               = MemState_IsPoolAllocated // < Alias for \ref MemState_IsPoolAllocated.
	MemState_MapProcessAllowed          = 1 << 23                  // < Map process allowed.
	MemState_AttrChangeAllowed          = 1 << 24                  // < Attribute change allowed.
	MemState_CodeMemAllowed             = 1 << 25                  // < Code memory allowed.
)

/// Memory attribute bitmasks.
const (
	MemAttr_IsBorrowed     = 1 << 0 // < Is borrowed memory.
	MemAttr_IsIpcMapped    = 1 << 1 // < Is IPC mapped (when IpcRefCount > 0).
	MemAttr_IsDeviceMapped = 1 << 2 // < Is device mapped (when DeviceRefCount > 0).
	MemAttr_IsUncached     = 1 << 3 // < Is uncached.
)

// Memory permission bitmasks.
const (
	Perm_None     Permission = 0               // < No permissions.
	Perm_R        Permission = 1 << 0          // < Read permission.
	Perm_W        Permission = 1 << 1          // < Write permission.
	Perm_X        Permission = 1 << 2          // < Execute permission.
	Perm_Rw       Permission = Perm_R | Perm_W // < Read/write permissions.
	Perm_Rx       Permission = Perm_R | Perm_X // < Read/execute permissions.
	Perm_DontCare Permission = 1 << 28         // < Don't care
)
