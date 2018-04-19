package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"timebank/src/cli"
	"timebank/src/migration"
	"timebank/src/operation"
)

func main() {
	defer recoverHandler()

	flags := cli.InitFlags()

	if flags.InterativeMode {
		cli.Console(flags.Filename)
	} else {
		collectionGetter := migration.Collections(flags.Filename)
		saver, saverCloser := operation.MakeSaver()

		defer saverCloser()

		collections := strings.Split(flags.Collections, ",")
		regex := regexp.MustCompile(`(\w+)(?:\[(\d+)?:(\d+)\])?`)

		for _, name := range collections {
			match := regex.FindStringSubmatch(name)
			name := match[1]
			start, _ := strconv.Atoi(match[2])
			end, _ := strconv.Atoi(match[3])

			collection := collectionGetter(name)

			if flags.SaveOperation {
				saver(collection, name)
			} else {
				operation.Print(collection, start, end)
			}
		}
	}
}

func recoverHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
		cli.PrintDefaults()
	}
}
