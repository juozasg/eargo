package main

import (
	"fmt"

	"github.com/juozasg/portmidi"
)

var midiInStream, midiOutStream *portmidi.Stream

var midiEvents = make(<-chan portmidi.Event)
var noteEvents = make(chan portmidi.Event)

const (
	midiNoteOn  = 0x90
	midiNoteOff = 0x80
)

func filterNoteEvents() {
	for {
		e := <-midiEvents
		if e.Status == midiNoteOn || e.Status == midiNoteOff {
			noteEvents <- e
		}
	}
}

func connectMIDIInput() {
	// select the default input device
	fmt.Println("Using the default MIDI input device")
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
		quit <- 1
		return
	}

	midiEvents = midiInStream.Listen()
	deviceInfo := portmidi.Info(deviceID)
	fmt.Printf("Listening to MIDI input device: #%d %s (%s)\n",
		deviceID, deviceInfo.Name, deviceInfo.Interface)

	go filterNoteEvents()
}

func connectMIDIOutput() {
	fmt.Println("")
	deviceId := -1
	for id := 0; id < portmidi.CountDevices(); id++ {
		info := portmidi.Info(portmidi.DeviceID(id))
		if info.Name == "eargo-fluidsynth" {
			deviceId = id
		}
	}

	if deviceId == -1 {
		fmt.Println("No MIDI output device named 'eargo-fluidsynth' found")
		quit <- 1
		return
	}

	var err error
	midiOutStream, err = portmidi.NewOutputStream(portmidi.DeviceID(deviceId), 1024, 0)
	if err != nil {
		fmt.Println("MIDI ERROR", err)
		quit <- 1
		return
	}
}

func connectMIDI() {
	fmt.Println("Initializing MIDI...")
	portmidi.Initialize()

	connectMIDIInput()
	connectMIDIOutput()
}

func cleanupMIDI() {
	fmt.Println("Closing MIDI devices...")
	if midiInStream != nil {
		midiInStream.Close()
	}

	if midiOutStream != nil {
		midiOutStream.Close()
	}

	portmidi.Terminate()
}
