package main

import (
	"fmt"

	"github.com/rakyll/portmidi"
)

var midiInStream *portmidi.Stream
var midiEvents = make(<-chan portmidi.Event)

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

	deviceInfo := portmidi.Info(deviceID)
	fmt.Printf("Using MIDI input device: #%d %s (%s)\n",
		deviceID, deviceInfo.Name, deviceInfo.Interface)

	var err error
	midiInStream, err = portmidi.NewInputStream(deviceID, 1024)
	if err != nil {
		fmt.Println("MIDI ERROR", err)
		return
	}

	midiEvents = midiInStream.Listen()
}

func cleanupMIDI() {
	if midiInStream != nil {
		midiInStream.Close()
	}
	portmidi.Terminate()
}
