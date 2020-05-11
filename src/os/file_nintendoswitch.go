// +build nintendoswitch

package os

import "unsafe"

// Read is unsupported on this system.
func (f *File) Read(b []byte) (n int, err error) {
	n = libc_read(int(f.fd), unsafe.Pointer(&b[0]), uint64(len(b)))

	return n, nil
}

// Write writes len(b) bytes to the output. It returns the number of bytes
// written or an error if this file is not stdout or stderr.
func (f *File) Write(b []byte) (n int, err error) {
	n = libc_write(int(f.fd), unsafe.Pointer(&b[0]), uint64(len(b)))

	return n, nil
}

// Close is unsupported on this system.
func (f *File) Close() error {
	return errUnsupported
}

// int write(int fd, const void *buf, size_t cnt)
//go:export write
func libc_write(fd int, buffer unsafe.Pointer, size uint64) int

// int read(int fd, void *buf, size_t count);
//go:export read
func libc_read(fd int, buffer unsafe.Pointer, size uint64) int
