OUTPUT_ARCH(aarch64)
ENTRY(_start)

PHDRS
{
  code   PT_LOAD FLAGS(5) /* Read | Execute */;
  rodata PT_LOAD FLAGS(4) /* Read */;
  data   PT_LOAD FLAGS(6) /* Read | Write */;
  dyn    PT_DYNAMIC;
}

_stack_size = 2K;

SECTIONS
{
  /* =========== CODE section =========== */
  PROVIDE(__start__ = 0x0);
  . = __start__;
  __code_start = . ;

  .crt0 :
  {
    KEEP (*(.crt0))
    . = ALIGN(8);
  } :code

  .init :
  {
    KEEP( *(.init) )
    . = ALIGN(8);
  } :code

  .plt :
  {
    *(.plt)
    *(.iplt)
    . = ALIGN(8);
  } :code

  .text :
  {
    *(.text.unlikely .text.*_unlikely .text.unlikely.*)
    *(.text.exit .text.exit.*)
    *(.text.startup .text.startup.*)
    *(.text.hot .text.hot.*)
    *(.text .stub .text.* .gnu.linkonce.t.*)
    . = ALIGN(8);
  } :code

  .fini :
  {
    KEEP( *(.fini) )
    . = ALIGN(8);
  } :code

  /* =========== RODATA section =========== */
  . = ALIGN(0x1000);
  __rodata_start = . ;

  .nx-module-name : { KEEP (*(.nx-module-name)) } :rodata

  .rodata :
  {
    *(.rodata .rodata.* .gnu.linkonce.r.*)
    . = ALIGN(8);
  } :rodata

  .eh_frame_hdr     : { __eh_frame_hdr_start = .; *(.eh_frame_hdr) *(.eh_frame_entry .eh_frame_entry.*) __eh_frame_hdr_end = .; } :rodata
  .eh_frame         : ONLY_IF_RO { KEEP (*(.eh_frame)) *(.eh_frame.*) } :rodata

  .dynamic           : { *(.dynamic) } :rodata :dyn
  .dynsym            : { *(.dynsym) } :rodata
  .dynstr            : { *(.dynstr) } :rodata
  .rela.dyn          : { *(.rela.*) } :rodata
  .interp            : { *(.interp) } :rodata


  /* Start address (in flash) of .data, used by startup code. */
  _sidata = LOADADDR(.data);
  /* =========== DATA section =========== */
  . = ALIGN(0x1000);
  __data_start = . ;

  __got_start__ = .;

  .got            : { *(.got) *(.igot) } :data
  .got.plt        : { *(.got.plt)  *(.igot.plt) } :data

  __got_end__ = .;

  .data ALIGN(8) :
  {
    _sdata = .;        /* used by startup code */
    *(.data .data.* .gnu.linkonce.d.*)
    SORT(CONSTRUCTORS)
    _edata = .;        /* used by startup code */
  } :data

  __bss_start__ = .;
  .bss ALIGN(8) :
  {
    _sbss = .;         /* used by startup code */
    *(.dynbss)
    *(.bss .bss.* .gnu.linkonce.b.*)
    *(COMMON)
    . = ALIGN(8);
    _ebss = .;         /* used by startup code */
  } : data
  __bss_end__ = .;

  __end__ = ABSOLUTE(.) ;

  . = ALIGN(0x1000);
  __argdata__ = ABSOLUTE(.) ;

  /* ==================
     ==== Metadata ====
     ================== */

  /* Discard sections that difficult post-processing */
  /DISCARD/ : { *(.group .comment .note) }
}

/* For the memory allocator. */
_globals_start = _sdata;
_globals_end = _ebss;
