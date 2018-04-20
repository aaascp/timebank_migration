package utils

import (
	"bufio"
	"fmt"
	"os"
)

func Confirm(message string) bool {
	input := make(chan string)
	reader := MakeReader(input)

	fmt.Println(message, " (y/n)")
	go reader()

	answer := <-input
	switch answer {
	case "y":
		return true
	default:
		return false
	}
}

func MakeReader(userInput chan string) func() {
	input := userInput

	return func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")

		text, _ := reader.ReadString('\n')
		input <- text
	}
}
