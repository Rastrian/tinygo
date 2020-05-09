// +build nintendoswitch

// +build gc.conservative gc.leaking

package runtime

var heapStart, heapEnd uintptr

func preinit() {
	size := uint64(0x2000000 * 16)
	memTotal, _ := SvcGetInfo(InfoType_TotalMemorySize, CUR_PROCESS_HANDLE, 0)
	memUsed, _ := SvcGetInfo(InfoType_UsedMemorySize, CUR_PROCESS_HANDLE, 0)

	if memTotal > memUsed+0x200000 {
		size = (memTotal - memUsed - 0x200000) & ^uint64(0x1FFFFF)
	}

	outAddr, result := SvcSetHeapSize(size)

	if result != ResultOK {
		panic("Cannot allocate heap")
	}

	heapStart = outAddr
	heapEnd = heapStart + uintptr(size)
}
