package main

import (
	"fmt"
	"machine"
)

func main() {
	machine.ConsoleInit(nil)
	fmt.Println("Golang PRINT!!!")

	for machine.AppletMainLoop() {
		machine.HidScanInput()

		keysDown := machine.HidKeysDown(10)

		if keysDown&(1<<10) > 0 {
			break
		}

		fmt.Printf("GoLang KeysDown: %032b\n", keysDown)

		machine.ConsoleUpdate(nil)
	}

	machine.ConsoleExit(nil)
}
