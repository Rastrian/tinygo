.section ".crt0","ax"
.global _start
.align

_start:
    b startup
    .word __nx_mod0 - _start
    .ascii "HOMEBREW"

.org _start+0x80
startup:
    // save lr
    mov  x7, x30

    // get aslr base
    bl   +4
    sub  x6, x30, #0x88

    b bssclr_start

bssclr_start:
    mov x27, x7
    mov x25, x5
    mov x26, x4

    // clear .bss
    adrp x0, __bss_start__
    adrp x1, __bss_end__
    add  x0, x0, #:lo12:__bss_start__
    add  x1, x1, #:lo12:__bss_end__
    sub  x1, x1, x0  // calculate size
    add  x1, x1, #7  // round up to 8
    bic  x1, x1, #7

bss_loop:
    str  xzr, [x0], #8
    subs x1, x1, #8
    bne  bss_loop

    b main

.global __nx_exit
.type   __nx_exit, %function
__nx_exit:
    // jump back to loader
    br   x1

.global __nx_mod0
__nx_mod0:
    .ascii "MOD0"
    .word  0
    .word  __bss_start__        - __nx_mod0
    .word  __bss_end__          - __nx_mod0
    .word  __eh_frame_hdr_start - __nx_mod0
    .word  __eh_frame_hdr_end   - __nx_mod0
    .word  0 // "offset to runtime-generated module object" (neither needed, used nor supported in homebrew)

    // MOD0 extensions for homebrew
    .ascii "LNY0"
    .word  __got_start__        - __nx_mod0
    .word  __got_end__          - __nx_mod0
