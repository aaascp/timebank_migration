package collection

import (
	"fmt"
	"strings"
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
	dict["subcategories"] = category.Subcategories
	return dict
}
