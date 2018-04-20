package cli

import (
	"bufio"
	"fmt"
	"os"
)

func Confirm(message string) bool {
	input := make(chan string)
	reader := makeReader(input)

	fmt.Println(message, " (y/n)")
	go reader()

	answer := <-input
	switch answer {
	case "y\n":
		return true
	default:
		return false
	}
}

func makeReader(userInput chan string) func() {
	input := userInput

	return func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")

		text, _ := reader.ReadString('\n')
		input <- text
	}
}
