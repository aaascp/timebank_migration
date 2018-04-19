package operation

import (
	"fmt"
	"timebank/src/collection"
	"timebank/src/migration"
)

type Printer func(string, int, int)

func MakePrinter(filename string) Printer {
	var items map[string][]collection.Collection

	return func(name string, start, end int) {
		if list := items[name]; list == nil {
			fmt.Printf("***** Fetching: %s *****\n", name)
			items = migration.Collections(filename, name)
		}

		actualCollection := items[name]
		fmt.Printf("***** %s *****\n", name)
		if end == 0 {
			end = len(actualCollection)
		}

		for i, item := range actualCollection[start:end] {
			fmt.Printf("%d - %s\n", i+1, item.ToString())
		}
	}
}
