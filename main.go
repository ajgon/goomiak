package main

import (
	"z80/machine"
)

func main() {
	machine := machine.NewSpectrum48k()
	machine.LoadFileToMemory(0x0000, "./roms/48.rom")
	machine.Run()
}
