package operation

import (
	"fmt"
	"log"
	"timebank/src/collection"

	mgo "gopkg.in/mgo.v2"
)

type Closer func()
type Saver func(string)

func MakeSaver(filename string) (Saver, Closer) {
	var items map[string][]collection.Collection

	session, err := mgo.Dial("mongodb://admin:admin@ds143039.mlab.com:43039/timebank")
	if err != nil {
		panic(err)
	}

	closeSession := func() {
		session.Close()
	}

	save := func(collectionName string) {
		collection := session.DB("timebank").C(collectionName)

		if list := items[collectionName]; list == nil {
			items = possibleCollections(filename, collectionName)
		}

		actualCollection := items[collectionName]
		persistableCollection := make([]interface{}, len(actualCollection))

		for i, item := range items[collectionName] {
			if item == nil {
				fmt.Println("nil")
			}
			persistableCollection[i] = item.ToDbFormat()
		}

		collection.DropCollection()
		err = collection.Insert(persistableCollection...)
		if err != nil {
			log.Fatal(err)
		}
	}

	return save, closeSession
}
