package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	mt "github.com/brettbuddin/musictheory"
)

var fluidsynthCmd *exec.Cmd

const fluidsynthName = "eargo-fluidsynth"

func startFluidsynth() {
	fluidsynthCmd = exec.Command("fluidsynth", "soundfont.sf2", "-s", "-i", "-p", fluidsynthName)
	fluidsynthCmd.Stdout = os.Stdout
	fluidsynthCmd.Stderr = os.Stderr

	fmt.Println("Starting fluidsynth with:", fluidsynthCmd)
	if err := fluidsynthCmd.Start(); err != nil {
		fmt.Println("Command error", fluidsynthCmd, err)
		quit <- 1
		return
	}

	// give it some time to start before portmidi initializes
	time.Sleep(2000 * time.Millisecond)
	fmt.Println("Started fluidsynth")
}

func stopFluidsynth() {
	if fluidsynthCmd != nil {
		// Kill it:
		if err := fluidsynthCmd.Process.Kill(); err != nil {
			fmt.Println("Failed to kill process: ", err)
		}

		fmt.Println("Stopped fluidsynth")
	}
}

func synthNote(event int64, p mt.Pitch) {
	// fmt.Println("bleeeoop ", p.Name(mt.AscNames))
	midiOutStream.WriteShort(event, int64(p.MIDI()), 127)
}
