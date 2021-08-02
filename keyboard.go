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
		mu.Lock()
		r, key, err := keyboard.GetKey()
		mu.Unlock()

		if err != nil {
			fmt.Println("Keyboard error:", err)
			return
		}

		input := int(key)
		if input == 0 {
			input = int(r)
		}

		keyboardInput <- input

		if key == keyboard.KeyEsc {
			quit <- 1
		}
	}
}

func keyboardReaderLoop() {
	fmt.Println("Using STDIN")

	fmt.Println("Type q to quit")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		r, _, _ := reader.ReadRune()
		if r == rune(113) {
			quit <- 1
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
	mu.Lock()

	keyboard.Close()

	mu.Unlock()
}
