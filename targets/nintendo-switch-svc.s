.macro SVC_BEGIN name
  .section .text.\name, "ax", %progbits
  .global \name
  .type \name, %function
  .align 2
  .cfi_startproc
\name:
.endm

.macro SVC_END
  .cfi_endproc
.endm

// func svcGetInfo(id0, handle uint32, id1 uint64) (value uint64, result Result)
// func svcGetInfo(w0, w1, x2) (x0, w1)
SVC_BEGIN svcGetInfo
  mov x3, x2
  mov w2, w1
  mov w1, w0
  svc 0x29
  mov x3, x1
  mov x1, x0
  mov x1, x3
  ret
SVC_END


// func svcSetHeapSize(size uint64) (outAddr uintptr, result Result)
// func svcSetHeapSize(x0) (x0 uintptr, w1 Result)
SVC_BEGIN svcSetHeapSize
    mov x1, x0
	svc 0x1
    mov x3, x1
    mov x1, x0
    mov x0, x3
	ret
SVC_END
