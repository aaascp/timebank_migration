package operation

import (
	"fmt"
	"timebank/src/collection"
	"timebank/src/migration"
)

type Items []interface{}

func possibleCollections(filename, name string) map[string][]collection.Collection {
	possibleCollections := make(map[string][]collection.Collection)

	if name == "users" || name == "transactions" {
		users, credits := migration.Users(filename)

		possibleCollections["users"] = users
		possibleCollections["transactions"] = credits

	} else if name == "services" || name == "categories" {
		services, categories := migration.Services(filename)

		possibleCollections["services"] = services
		possibleCollections["categories"] = categories
	} else {
		panic(fmt.Sprintf("Collection [%s] not exists", name))
	}

	return possibleCollections
}
