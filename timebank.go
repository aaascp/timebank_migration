package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"timebank/migration"

	mgo "gopkg.in/mgo.v2"
)

type Items []interface{}

var usageFlag = flag.Bool("u", false, "Usage")
var wbPath, _ = filepath.Abs("./resources/btf.xlsx")
var filename = flag.String("f", wbPath, "Path to BTF workbook")
var collections = flag.String(
	"c",
	"users[9:3961],credits[9:3961],services[4:6690],categories[4:6690]",
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

func possibleCollections(collection string) map[string][]interface{} {
	possibleCollections := make(map[string][]interface{})

	if collection == "users" || collection == "credits" {
		users, credits := migration.Users(*filename)

		possibleCollections["users"] = users
		possibleCollections["credits"] = credits

	} else if collection == "services" || collection == "categories" {
		services, categories := migration.Services(*filename)

		possibleCollections["services"] = services
		possibleCollections["categories"] = categories
	} else {
		panic(fmt.Sprintf("Collection [%s] not exists", collection))
	}

	return possibleCollections
}

func makePrinter() func(string, int, int) {
	var items map[string][]interface{}

	return func(collection string, start, end int) {
		if i := items[collection]; i != nil {
			fmt.Printf("***** Collection present: %s *****\n", collection)
			fmt.Println(i[start:end])
		} else {
			fmt.Printf("***** Collection NOT present: %s *****\n", collection)
			items = possibleCollections(collection)
			fmt.Println(items[collection][start:end])
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
