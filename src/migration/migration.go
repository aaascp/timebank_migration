package migration

import (
	"fmt"
	"timebank/src/collection"
)

type Items []interface{}

func Collections(filename string) func(string) []collection.Item {
	collections := make(map[string][]collection.Item)

	return func(name string) []collection.Item {
		if list := collections[name]; list == nil {
			fmt.Printf("***** Fetching: %s *****\n", name)
			migrate(filename, name, collections)
		}
		fmt.Printf("***** %s fetched *****\n", name)
		return collections[name]
	}
}

func migrate(filename, name string, collections map[string][]collection.Item) {
	if name == "user" || name == "transaction" {
		users, credits := Users(filename)

		collections["user"] = users
		collections["transaction"] = credits

	} else if name == "service" || name == "category" {
		services, categories := Services(filename)

		collections["service"] = services
		collections["category"] = categories
	} else {
		panic(fmt.Sprintf("Collection [%s] does not exists", name))
	}
}
