package main

import (
	"fmt"

	"github.com/juozasg/portmidi"
)

var midiInStream *portmidi.Stream

var midiEvents = make(<-chan portmidi.Event)
var noteEvents = make(chan portmidi.Event)

func filterNoteEvents() {
	for {
		e := <-midiEvents
		if e.Status == 0x80 || e.Status == 0x90 {
			noteEvents <- e
		}
	}
}

func startMIDILoop() {
	fmt.Println("Initializing MIDI...")
	portmidi.Initialize()

	// fmt.Println("Total MIDI Devices:", portmidi.CountDevices())

	deviceID := portmidi.DefaultInputDeviceID()
	if deviceID == -1 {
		fmt.Println("No default MIDI input device found")
		quit <- 1
		return
	}

	var err error
	midiInStream, err = portmidi.NewInputStream(deviceID, 1024)
	if err != nil {
		fmt.Println("MIDI ERROR", err)
		return
	}

	midiEvents = midiInStream.Listen()
	deviceInfo := portmidi.Info(deviceID)
	fmt.Printf("Listening to MIDI input device: #%d %s (%s)\n",
		deviceID, deviceInfo.Name, deviceInfo.Interface)

	go filterNoteEvents()
}

func cleanupMIDI() {
	fmt.Println("Closing MIDI device...")
	if midiInStream != nil {
		midiInStream.Close()
	}
	portmidi.Terminate()
}
