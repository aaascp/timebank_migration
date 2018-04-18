package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"timebank/src/collection"
	"timebank/src/migration"

	mgo "gopkg.in/mgo.v2"
)

type Items []interface{}

var usageFlag = flag.Bool("u", false, "Usage")
var wbPath, _ = filepath.Abs("./resources/btf.xlsx")
var filename = flag.String("f", wbPath, "Path to BTF workbook")
var collections = flag.String(
	"c",
	"users[:10],transactions[:10],services[:10],categories",
	"Collections separated by commas")

func main() {
	flag.Parse()
	if *usageFlag {
		flag.PrintDefaults()
		return
	}

	printer := makePrinter()
	collectionsList := strings.Split(*collections, ",")
	regex := regexp.MustCompile(`(\w+)(?:\[(\d+)?:(\d+)\])?`)
	for _, collection := range collectionsList {
		match := regex.FindStringSubmatch(collection)
		name := match[1]
		start, _ := strconv.Atoi(match[2])
		end, _ := strconv.Atoi(match[3])
		printer(name, start, end)
	}
}

func possibleCollections(name string) map[string][]collection.Collection {
	possibleCollections := make(map[string][]collection.Collection)

	if name == "users" || name == "transactions" {
		// users, credits := migration.Users(*filename)
		//
		// possibleCollections["users"] = users
		// possibleCollections["transactions"] = credits

	} else if name == "services" || name == "categories" {
		services, categories := migration.Services(*filename)

		possibleCollections["services"] = services
		possibleCollections["categories"] = categories
	} else {
		panic(fmt.Sprintf("Collection [%s] not exists", name))
	}

	return possibleCollections
}

func makePrinter() func(string, int, int) {
	var items map[string][]collection.Collection

	return func(name string, start, end int) {
		if list := items[name]; list == nil {
			fmt.Printf("***** Fetching: %s *****\n", name)
			items = possibleCollections(name)
		}

		fmt.Printf("***** %s *****\n", name)
		if end == 0 {
			end = len(items[name])
		}

		for i, item := range items[name][start:end] {
			fmt.Printf("%d - %s\n", i+1, item.ToString())
		}
	}
}

func makeSaver() func(Items, string) {
	session, err := mgo.Dial("mongodb://admin:admin@ds143039.mlab.com:43039/timebank")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	return func(items Items, collectionName string) {
		collection := session.DB("timebank").C(collectionName)

		collection.DropCollection()
		err = collection.Insert(items...)
		if err != nil {
			log.Fatal(err)
		}
	}
}
