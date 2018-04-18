package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"timebank/src/cli"
	"timebank/src/operation"
)

func main() {
	defer recoverHandler()

	flags := cli.InitFlags()
	printer := operation.MakePrinter(flags.Filename)
	saver, saverCloser := operation.MakeSaver(flags.Filename)

	defer saverCloser()

	collections := strings.Split(flags.Collections, ",")
	regex := regexp.MustCompile(`(\w+)(?:\[(\d+)?:(\d+)\])?`)

	for _, collection := range collections {
		match := regex.FindStringSubmatch(collection)
		name := match[1]
		start, _ := strconv.Atoi(match[2])
		end, _ := strconv.Atoi(match[3])

		if flags.SaveOperation {
			saver(name)
		} else {
			printer(name, start, end)
		}
	}
}

func recoverHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
		cli.PrintDefaults()
	}
}
