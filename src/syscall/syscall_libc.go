// +build darwin nintendoswitch

package syscall

import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	notImplemented = errors.New("os: not implemented")
)

type Signal int

func Close(fd int) (err error) {
	errN := libc_close(int32(fd))
	if errN != 0 {
		return fmt.Errorf("error %d", errN)
	}
	return nil
}

// getErrno returns the current C errno. It may not have been caused by the last
// call, so it should only be relied upon when the last call indicates an error
// (for example, by returning -1).
func getErrno() Errno {
	errptr := libc__errno()
	return Errno(errptr)
}

func Write(fd int, p []byte) (n int, err error) {
	buf, count := splitSlice(p)
	n = libc_write(int32(fd), buf, uint(count))
	if n < 0 {
		err = getErrno()
	}
	return
}

func Read(fd int, p []byte) (n int, err error) {
	buf, count := splitSlice(p)
	n = libc_read(int32(fd), buf, uint(count))
	if n < 0 {
		err = getErrno()
	}
	return
}

func Seek(fd int, offset int64, whence int) (off int64, err error) {
	off = libc_lseek(int32(fd), offset, whence)

	return off, nil
}

func Open(path string, mode int, perm uint32) (fd int, err error) {
	buf, _ := splitString(path)
	fd = libc_open(buf, uint(mode), uint(perm))
	if fd < 0 {
		err = getErrno()
	}

	return
}

func Mkdir(path string, mode uint32) (err error) {
	return ENOSYS // TODO
}

func Unlink(path string) (err error) {
	return ENOSYS // TODO
}

func Kill(pid int, sig Signal) (err error) {
	return notImplemented // TODO
}

func Getpid() (pid int) {
	panic("unimplemented: getpid") // TODO
}

func Getenv(key string) (value string, found bool) {
	return "", false // TODO
}

func splitSlice(p []byte) (buf *byte, len uintptr) {
	slice := (*struct {
		buf *byte
		len uintptr
		cap uintptr
	})(unsafe.Pointer(&p))
	return slice.buf, slice.len
}

func splitString(p string) (buf *byte, len uintptr) {
	slice := (*struct {
		ptr    *byte
		length uintptr
	})(unsafe.Pointer(&p))
	return slice.ptr, slice.length
}

// int write(int fd, const void *buf, size_t cnt)
//go:export write
func libc_write(fd int32, buffer *byte, size uint) int

// int read(int fd, void *buf, size_t count);
//go:export read
func libc_read(fd int32, buffer *byte, size uint) int

// int close(int fd);
//go:export read
func libc_close(fd int32) int

// int open(const char *pathname, int flags, mode_t mode);
//go:export open
func libc_open(pathname *byte, flags uint, mode uint) int

//go:export lseek
func libc_lseek(fd int32, offset int64, whence int) int64

//go:export __errno
func libc__errno() uintptr
