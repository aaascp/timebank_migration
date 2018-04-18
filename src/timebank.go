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
var allFlag = flag.Bool("a", false, "If this flag is set all collections will be processed")
var saveFlag = flag.Bool("s", false, "If this flag is set the collection will be saved to database")

var wbPath, _ = filepath.Abs("./resources/btf.xlsx")
var filename = flag.String("f", wbPath, "Path to BTF workbook")
var collections = flag.String(
	"c",
	"",
	"Collections separated by commas")

func main() {
	defer recoverHandler()

	initFlags()

	printer := makePrinter()
	saverCloser, saver := makeSaver()

	defer saverCloser()

	collectionsList := strings.Split(*collections, ",")
	regex := regexp.MustCompile(`(\w+)(?:\[(\d+)?:(\d+)\])?`)
	for _, collection := range collectionsList {
		match := regex.FindStringSubmatch(collection)
		name := match[1]
		start, _ := strconv.Atoi(match[2])
		end, _ := strconv.Atoi(match[3])

		if *saveFlag {
			saver(name)
		} else {
			printer(name, start, end)
		}
	}
}

func possibleCollections(name string) map[string][]collection.Collection {
	possibleCollections := make(map[string][]collection.Collection)

	if name == "users" || name == "transactions" {
		users, credits := migration.Users(*filename)

		possibleCollections["users"] = users
		possibleCollections["transactions"] = credits

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

func makeSaver() (func(), func(string)) {
	var items map[string][]collection.Collection

	session, err := mgo.Dial("mongodb://admin:admin@ds143039.mlab.com:43039/timebank")
	if err != nil {
		panic(err)
	}

	closeSession := func() {
		session.Close()
	}

	return closeSession, func(collectionName string) {
		collection := session.DB("timebank").C(collectionName)

		if list := items[collectionName]; list == nil {
			items = possibleCollections(collectionName)
		}

		actualCollection := items[collectionName]
		persistableCollection := make([]interface{}, len(actualCollection))

		for i, item := range actualCollection {
			if item.ToDbFormat() == nil {
				fmt.Println(i, item, "Nil!")
			}
			persistableCollection = append(persistableCollection, item.ToDbFormat())
		}

		collection.DropCollection()
		err = collection.Insert(persistableCollection...)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func initFlags() {
	flag.Parse()

	if *usageFlag {
		panic("Usage:")
	}

	if *allFlag {
		*collections = "users[:10],transactions[:10],services[:10],categories"
	}

	if !*allFlag && *collections == "" {
		panic("Inform collections of interest with -c or -a flags.\nUsage:")
	}
}

func recoverHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
		flag.PrintDefaults()
	}
}
