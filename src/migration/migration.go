package migration

import (
	"fmt"
	"timebank/src/collection"
)

type Items []interface{}

func Collections(filename, name string) map[string][]collection.Collection {
	collections := make(map[string][]collection.Collection)

	if name == "user" || name == "transaction" {
		users, credits := Users(filename)

		collections["user"] = users
		collections["transaction"] = credits

	} else if name == "service" || name == "category" {
		services, categories := Services(filename)

		collections["service"] = services
		collections["category"] = categories
	} else {
		panic(fmt.Sprintf("Collection [%s] not exists", name))
	}

	return collections
}
