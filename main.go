package main

import (
	"flag"
	"z80/loader"
	"z80/machine"
)

var z80FileFlag = flag.String("z80file", "", "path to z80 file to be loaded")
var tapFileFlag = flag.String("tapfile", "", "path to tap file to be loaded")
var romFileFlag = flag.String("romfile", "roms/48.rom", "path to rom file to be used")
var fullSpeedFlag = flag.Bool("fullspeed", false, "run emulator with full speed")

func main() {
	flag.Parse()

	machine := machine.NewSpectrum48k()
	machine.LoadFileToMemory(0x0000, *romFileFlag)
	machine.FullSpeed(*fullSpeedFlag)

	if *tapFileFlag != "" {
		tapFile := loader.NewTapFile(*tapFileFlag)
		machine.InsertTape(tapFile)
	} else if *z80FileFlag != "" {
		snapshot := loader.Z80(*z80FileFlag)
		machine.LoadSnapshot(snapshot)
	}

	machine.Run()
}
