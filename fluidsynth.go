package main

import (
	"fmt"
	// "time"

	mt "github.com/brettbuddin/musictheory"
)

func startFluidsynth() {
	fmt.Println("Started fluidsynth")
	//bin/fluidsynth share/fluid-synth/sf2/VintageDreamsWaves-v2.sf2  -p eargo-fluidsynth -s -i
}

func synthNote(event int64, p mt.Pitch) {
	// fmt.Println("bleeeoop ", p.Name(mt.AscNames))
	midiOutStream.WriteShort(event, int64(p.MIDI()), 127)
}
