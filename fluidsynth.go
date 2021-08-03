package main

import (
	"fmt"
	mt "github.com/brettbuddin/musictheory"
)

func startFluidsynth() {
	fmt.Println("Started fluidsynth")
}

func playNote(p mt.Pitch) {
	fmt.Println("bleeeoop ", p.Name(mt.AscNames))
}
