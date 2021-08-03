package main

import (
	"fmt"

	mt "github.com/brettbuddin/musictheory"
	"github.com/brettbuddin/musictheory/intervals"
)

/*
	TODO:
	1) read midi in
	2) determine the note pitch
	0) play a sound with a pitch
*/

func main() {
	fmt.Println("music theory data structures:")
	root := mt.NewPitch(mt.C, mt.Natural, 4)
	fmt.Println(root.Name(mt.AscNames), "MIDI", root.MIDI())
	fmt.Println(mt.NewScale(root, intervals.Phrygian, 1))

	startKeyoardIOLoop()
	defer cleanupKeyboard()

	go startMIDILoop()
	defer cleanupMIDI()

	for {
		select {
		case <-quit:
			fmt.Println("Exiting!")
			return
		case <-keyboardInput:
			// fmt.Printf("%d", i)
		case e := <-midiEvents:
			fmt.Println(e)
		}
	}
}
