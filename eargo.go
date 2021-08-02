package main

import (
	// "github.com/brettbuddin/musictheory/intervals"
	"fmt"

	mt "github.com/brettbuddin/musictheory"
)

/*
	TODO:
	1) read midi in
	2) determine the note pitch
	0) play a sound with a pitch
*/

func main() {
	root := mt.NewPitch(mt.C, mt.Natural, 4)
	fmt.Println("hello!")
	fmt.Println(root.Name(mt.AscNames), "MIDI", root.MIDI())
}
