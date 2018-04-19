package operation

import (
	"fmt"
	"timebank/src/collection"
)

func Print(items []collection.Item, start, end int) {
	if end == 0 {
		end = len(items)
	}

	for i, item := range items[start:end] {
		fmt.Printf("%d - %s\n", i+1, item.ToString())
	}
}
