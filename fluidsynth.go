package main

import (
	"fmt"
	"time"

	mt "github.com/brettbuddin/musictheory"
)

func startFluidsynth() {
	fmt.Println("Started fluidsynth")
}

func playNote(p mt.Pitch) {
	// fmt.Println("bleeeoop ", p.Name(mt.AscNames))

	midiOutStream.WriteShort(midiNoteOn, int64(p.MIDI()), 100)
	time.Sleep(time.Millisecond * 200)
	midiOutStream.WriteShort(midiNoteOff, int64(p.MIDI()), 100)
}
