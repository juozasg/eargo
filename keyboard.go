package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

var quit = make(chan int)
var keyboardInput = make(chan int)

func keyboardTTYLoop() {
	fmt.Println("** Press ESC to quit **")
	for {
		r, key, err := keyboard.GetKey()

		if err != nil {
			// ignore keyboard.Close() error
			if err.Error() != "operation canceled" {
				fmt.Println("Keyboard error:", err)
			}
			return
		}

		input := int(key)
		if input == 0 {
			input = int(r)
		}

		keyboardInput <- input

		if key == keyboard.KeyEsc {
			quit <- 1
			return
		}
	}
}

func keyboardReaderLoop() {
	fmt.Println("Type q to quit")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		r, _, _ := reader.ReadRune()
		if r == rune(113) {
			quit <- 1
			return
		}

		keyboardInput <- int(r)
	}
}

func startKeyoardIOLoop() {
	if err := keyboard.Open(); err == nil {
		go keyboardTTYLoop()
	} else {
		go keyboardReaderLoop()
	}
}

func cleanupKeyboard() {
	fmt.Println("Keyboard cleanup")
	keyboard.Close()
}
