package collection

import (
	"fmt"
	"strings"

	mgo "gopkg.in/mgo.v2"
)

type Category struct {
	Name          string
	Subcategories []string
}

func (category Category) ToString() string {
	return fmt.Sprintf(
		"{name: %s,subcategories: [%s]}",
		category.Name,
		strings.Join(category.Subcategories, ", "))
}

func (category Category) ToDbFormat() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["name"] = category.Name
	if len(category.Subcategories) != 0 {
		dict["subcategories"] = category.Subcategories
	}
	return dict
}

func CategoryIndexes() []mgo.Index {
	indexes := make([]mgo.Index, 2)

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"name"},
			Unique: true})

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"subcategories.name"},
			Unique: true,
			Sparse: true})

	return indexes
}
