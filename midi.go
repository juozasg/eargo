package main

import (
	"fmt"
	"github.com/rakyll/portmidi"
)

var midiInStream *portmidi.Stream
var midiEvents = make(<-chan portmidi.Event)

func startMIDILoop() {
	portmidi.Initialize()

	fmt.Println("Total MIDI Devices:", portmidi.CountDevices())

	deviceID := portmidi.DefaultInputDeviceID()
	deviceInfo := portmidi.Info(deviceID)
	fmt.Printf("Using the default MIDI input device: #%d\n", deviceID)
	fmt.Println("Device info: ", deviceInfo.Interface, deviceInfo.Name)

	var err error
	midiInStream, err = portmidi.NewInputStream(deviceID, 1024)
	if err != nil {
		fmt.Println("MIDI ERROR", err)
		return
	}

	midiEvents = midiInStream.Listen()
}

func cleanupMIDI() {
	mu.Lock()

	if midiInStream != nil {
		midiInStream.Close()
	}
	portmidi.Terminate()

	mu.Unlock()
}
