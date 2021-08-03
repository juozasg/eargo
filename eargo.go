package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mt "github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/musictheory/intervals"
)

func cleanup() {
	cleanupKeyboard()
	cleanupMIDI()
	fmt.Println("Bye")
	os.Exit(0)
}

var sigchan = make(chan os.Signal, 1)

func prepareCleanup() {
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigchan:
			cleanup()
		case <-quit:
			cleanup()
		}
	}()
}

func main() {
	fmt.Println("music theory data structures:")
	root := mt.NewPitch(mt.C, mt.Natural, 4)
	fmt.Println(root.Name(mt.AscNames), "MIDI", root.MIDI())
	fmt.Println(mt.NewScale(root, intervals.Phrygian, 1))

	prepareCleanup()
	startKeyoardIOLoop()
	connectMIDI()
	startFluidsynth()

	for {
		select {
		case i := <-keyboardInput:
			fmt.Printf("Key: %d\n", i)
		case e := <-noteEvents:
			semitones := int(e.Data1) - 24
			pitch := mt.Pitch{Interval: mt.Semitones(semitones)}
			fmt.Printf("Note %#x: %d\n", e.Status, e.Data1)

			if e.Status == midiNoteOn {
				playNote(pitch)
			}
		}
	}
}
