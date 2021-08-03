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

func dumpMIDIDevices() {
	fmt.Println("Available MIDI devices:")
	for id := 0; id < portmidi.CountDevices(); id++ {
		info := portmidi.Info(portmidi.DeviceID(id))
		fmt.Println(id, "-", info)
	}
	fmt.Println("")
}

func connectMIDIInput() {
	// select the default input device
	deviceId := portmidi.DefaultInputDeviceID()
	if deviceId == -1 {
		fmt.Println("No default MIDI input device found")
		quit <- 1
		return
	}

	var err error
	midiInStream, err = portmidi.NewInputStream(deviceId, 1024)
	if err != nil {
		fmt.Println("MIDI Error", err)
		quit <- 1
		return
	}

	midiEvents = midiInStream.Listen()
	deviceInfo := portmidi.Info(deviceId)
	fmt.Printf("Listening to the default MIDI input device: %d - %s (%s)\n",
		deviceId, deviceInfo.Name, deviceInfo.Interface)

	go filterNoteEvents()
}

// func waitForFluidsynth() {
// 	deviceId, _ := fluidsynthDeviceId()

// 	for retries := 0; retries < 200 && deviceId == -1; retries++ {
// 		time.Sleep(time.Millisecond * 1000)
// 		portmidi.Terminate()
// 		portmidi.Initialize()

// 		deviceId, _ = fluidsynthDeviceId()

// 		if deviceId != -1 {
// 			return
// 		}
// 	}
// }

func fluidsynthDeviceId() (int, *portmidi.DeviceInfo) {
	for id := 0; id < portmidi.CountDevices(); id++ {
		info := portmidi.Info(portmidi.DeviceID(id))
		if info.Name == fluidsynthName {
			return id, info
		}
	}

	return -1, nil
}

func connectMIDIOutput() {
	deviceId, deviceInfo := fluidsynthDeviceId()

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

	fmt.Printf("Sending to MIDI output device: #%d %s (%s)\n",
		deviceId, deviceInfo.Name, deviceInfo.Interface)
}

func connectMIDI() {
	fmt.Println("Initializing MIDI...")
	portmidi.Initialize()

	// waitForFluidsynth()
	dumpMIDIDevices()

	connectMIDIOutput()
	connectMIDIInput()
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
