// +build nintendoswitch

package machine

import "unsafe"

//go:export consoleInit
func ConsoleInit(c unsafe.Pointer)

//go:export appletMainLoop
func AppletMainLoop() bool

//go:export consoleUpdate
func ConsoleUpdate(c unsafe.Pointer)

//go:export consoleExit
func ConsoleExit(c unsafe.Pointer)

//go:export hidScanInput
func HidScanInput()

//u64 hidKeysDown(HidControllerID id);
//go:export hidKeysDown
func HidKeysDown(id uint64) uint64
