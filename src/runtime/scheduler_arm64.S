.section .text.tinygo_scanCurrentStack
.global  tinygo_scanCurrentStack
.type    tinygo_scanCurrentStack, %function
tinygo_scanCurrentStack:
   // x9 to x15: used to hold local variables (caller saved)
   // Push callee-saved registers onto the stack.
   sub  sp, sp, 256     // Open up some temp stack space
                        // Each X register has 64 bits, so 8 bytes
                        // Each pair then has 16 bytes
                        // There is 31 registers (0 to 30)
                        // So 31 * 8 = 248 bytes on the stack

   stp  x0, x1,   [sp, #(16 * 0 )]  // save 2 pairs of registers
   stp  x2, x3,   [sp, #(16 * 1 )]  // save 2 pairs of registers
   stp  x4, x5,   [sp, #(16 * 2 )]  // save 2 pairs of registers
   stp  x6, x7,   [sp, #(16 * 3 )]  // save 2 pairs of registers
   stp  x8, x9,   [sp, #(16 * 4 )]  // save 2 pairs of registers
   stp  x10, x11, [sp, #(16 * 5 )]  // save 2 pairs of registers
   stp  x12, x13, [sp, #(16 * 6 )]  // save 2 pairs of registers
   stp  x14, x15, [sp, #(16 * 7 )]  // save 2 pairs of registers
   stp  x16, x17, [sp, #(16 * 8 )]  // save 2 pairs of registers
   stp  x18, x19, [sp, #(16 * 9 )]  // save 2 pairs of registers
   stp  x20, x21, [sp, #(16 * 10)]  // save 2 pairs of registers
   stp  x22, x23, [sp, #(16 * 11)]  // save 2 pairs of registers
   stp  x24, x25, [sp, #(16 * 12)]  // save 2 pairs of registers
   stp  x26, x27, [sp, #(16 * 13)]  // save 2 pairs of registers
   stp  x28, x29, [sp, #(16 * 14)]  // save 2 pairs of registers
   stp  x30, x30, [sp, #(16 * 15)]  // save 2 pairs of registers

   // Scan the stack.
   mov x0, sp
   bl tinygo_scanstack

   // Restore stack state.
   add sp, sp, 256

   // Return to the caller.
   ret
