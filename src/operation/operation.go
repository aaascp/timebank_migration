package operation

import (
	"fmt"
	"timebank/src/collection"
	"timebank/src/migration"
)

type Items []interface{}

func possibleCollections(filename, name string) map[string][]collection.Collection {
	possibleCollections := make(map[string][]collection.Collection)

	if name == "user" || name == "transaction" {
		users, credits := migration.Users(filename)

		possibleCollections["user"] = users
		possibleCollections["transaction"] = credits

	} else if name == "service" || name == "category" {
		services, categories := migration.Services(filename)

		possibleCollections["service"] = services
		possibleCollections["category"] = categories
	} else {
		panic(fmt.Sprintf("Collection [%s] not exists", name))
	}

	return possibleCollections
}
