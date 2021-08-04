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

var keyboardTTYOpened = false

func startKeyoardIOLoop() {
	if err := keyboard.Open(); err == nil {
		keyboardTTYOpened = true
		fmt.Println("** Press ESC to quit **")
		go keyboardTTYLoop()
	} else {
		fmt.Println("Type q to quit")
		go keyboardReaderLoop()
	}
}

func cleanupKeyboard() {
	fmt.Println("Keyboard cleanup")
	if keyboardTTYOpened {
		keyboard.Close()
	}
}
