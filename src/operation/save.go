package operation

import (
	"timebank/src/collection"

	mgo "gopkg.in/mgo.v2"
)

type Saver func([]collection.Item, string)
type Closer func()

func MakeSaver() (Saver, Closer) {
	session, err := mgo.Dial("mongodb://admin:admin@ds143039.mlab.com:43039/timebank")
	if err != nil {
		panic(err)
	}

	closeSession := func() {
		session.Close()
	}

	save := func(items []collection.Item, name string) {
		dbCollection := session.DB("timebank").C(name)
		collections := persistableCollection(items)

		dbCollection.DropCollection()
		createIndex(name, dbCollection)

		err = dbCollection.Insert(collections...)
		if err != nil {
			panic(err)
		}
	}

	return save, closeSession
}

func createIndex(name string, dbCollection *mgo.Collection) {
	indexes := collection.Indexes[name]

	for _, index := range indexes {
		dbCollection.EnsureIndex(index)
	}
}

func persistableCollection(items []collection.Item) []interface{} {
	persistableCollection := make([]interface{}, len(items))

	for i, item := range items {
		persistableCollection[i] = item.ToDbFormat()
	}

	return persistableCollection
}
