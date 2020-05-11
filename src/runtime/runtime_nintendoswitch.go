// +build nintendoswitch

package runtime

//export usleep
func usleep(usec uint) int

//export abort
func abort() {
	exit(1)
}

//go:export exit
func exit(code int) int

//export clock_gettime
func clock_gettime(clk_id int32, ts *timespec)

type timeUnit int64

const tickMicros = 1

// Note: tv_sec and tv_nsec vary in size by platform. They are 32-bit on 32-bit
// systems and 64-bit on 64-bit systems (at least on macOS/Linux), so we can
// simply use the 'int' type which does the same.
type timespec struct {
	tv_sec  int // time_t: follows the platform bitness
	tv_nsec int // long: on Linux and macOS, follows the platform bitness
}

const CLOCK_MONOTONIC_RAW = 4

func postinit() {}

// Entry point for Go. Initialize all packages and call main.main().
//export main
func main() int {
	preinit()
	run()

	return exit(0) // Call libc_exit to cleanup libnx
}

func putchar(c byte) {
	OutputDebugChar(c)
}

const asyncScheduler = false

func sleepTicks(d timeUnit) {
	usleep(uint(d) / 1000)
}

// Return monotonic time in nanoseconds.
//
// TODO: noescape
func monotime() uint64 {
	ts := timespec{}
	clock_gettime(CLOCK_MONOTONIC_RAW, &ts)
	return uint64(ts.tv_sec)*1000*1000*1000 + uint64(ts.tv_nsec)
}

func ticks() timeUnit {
	return timeUnit(monotime())
}
