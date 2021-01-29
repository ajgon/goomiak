package main

import (
	"z80/machine"
)

func main() {
	machine := machine.NewSpectrum48k()
	machine.Run()
}
