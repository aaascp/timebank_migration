package collection

import mgo "gopkg.in/mgo.v2"

type Item interface {
	ToString() string
	ToDbFormat() map[string]interface{}
}

var Indexes = map[string][]mgo.Index{
	"user":        UserIndexes(),
	"transaction": TransactionIndexes(),
	"service":     ServiceIndexes(),
	"category":    CategoryIndexes(),
}
