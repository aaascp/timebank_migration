package collection

import (
	"fmt"
	"strings"
)

type Category struct {
	Name          string
	Subcategories []string
}

type ServiceCategory struct {
	Name        string
	Subcategory string
}

func (category ServiceCategory) ToString() string {
	return fmt.Sprintf(
		"{name: %s,subcategory: %s}",
		category.Name,
		category.Subcategory)
}

func (category Category) ToString() string {
	return fmt.Sprintf(
		"{name: %s,subcategories: [%s]}",
		category.Name,
		strings.Join(category.Subcategories, ", "))
}
