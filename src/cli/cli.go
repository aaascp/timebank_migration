package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"timebank/src/operation"
)

func Console(filename string) {
	input := make(chan string)
	isDone := make(chan bool)
	isDecoded := make(chan bool)

	decoder, decoderCloser := makeDecoder(filename, isDecoded, isDone)
	reader := makeReader(input)
	defer decoderCloser()

	go reader()
	for {

		select {
		case in := <-input:
			go decoder(in)
		case <-isDecoded:
			go reader()
		case <-isDone:
			fmt.Println("Bye.")
			return
		}
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

func makeDecoder(filename string, isDecoded, isDone chan bool) (func(string), func()) {
	printer := operation.MakePrinter(filename)
	saver, saverCloser := operation.MakeSaver(filename)

	regex := regexp.MustCompile(`(save|print|done)(?:\s+(\w+)(?:\[(\d+)?:(\d+)\])?)?`)

	decode := func(input string) {
		defer recoverHandler()

		match := regex.FindStringSubmatch(input)
		if len(match) < 4 {
			fmt.Println("Invalid operation")
		} else {
			action := match[1]
			collection := strings.TrimSpace(match[2])
			start, _ := strconv.Atoi(match[3])
			end, _ := strconv.Atoi(match[4])

			switch action {
			case "print":
				printer(collection, start, end)
			case "save":
				saver(collection)
				saverCloser()
			case "done":
				isDone <- true
			default:
				fmt.Println("Invalid operation")
			}
		}

		isDecoded <- true
	}

	return decode, saverCloser
}

func recoverHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}