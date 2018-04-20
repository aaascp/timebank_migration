package operation

import (
	"fmt"
	"log"
	"timebank/src/collection"
	"timebank/src/utils"

	mgo "gopkg.in/mgo.v2"
)

type Saver func([]collection.Item, string)
type Closer func()

func MakeSaver() (Saver, Closer) {
	confirmation := utils.Confirm("This operation will drop the current collection. Confirm?")
	if !confirmation {
		panic("Aborted")
	}

	session, err := mgo.Dial("mongodb://admin:admin@ds143039.mlab.com:43039/timebank")
	if err != nil {
		panic(err)
	}

	closeSession := func() {
		session.Close()
	}

	save := func(items []collection.Item, name string) {
		collection := session.DB("timebank").C(name)
		persistableCollection := make([]interface{}, len(items))

		for i, item := range items {
			persistableCollection[i] = item.ToDbFormat()
		}

		collection.DropCollection()
		err = collection.Insert(persistableCollection...)
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("Collection [%s] saved!\n", name)
		}
	}

	return save, closeSession
}
