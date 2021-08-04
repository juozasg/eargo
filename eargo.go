package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
)

var cleanedUp = false
var mu sync.Mutex

func cleanupPanic() {
	if r := recover(); r != nil {
		fmt.Println("PANIC", r)
		fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
	}

	cleanup()
}

func cleanup() {
	mu.Lock()

	if !cleanedUp {
		cleanupKeyboard()
		cleanupMIDI()
		stopFluidsynth()
		cleanedUp = true
	}
	mu.Unlock()

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
	// handle panics, os signals and quit channel signal
	defer cleanupPanic()
	prepareCleanup()

	startFluidsynth()
	connectMIDI()
	startKeyoardIOLoop()

	gameLoop()
}
