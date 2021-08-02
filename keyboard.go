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
	fmt.Println("Press ESC to quit")
	for {
		r, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", r, key)

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
		fmt.Println("read: ", r)
		if r == rune(113) {
			quit <- 1
		}

		keyboardInput <- int(r)
	}
}

func startKeyoardIOLoop() {
	if err := keyboard.Open(); err == nil {
		fmt.Println("keyboard open", err)
		go keyboardTTYLoop()
	} else {
		go keyboardReaderLoop()
	}
}
