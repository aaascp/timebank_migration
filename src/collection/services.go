package collection

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

type Service struct {
	UserName     string
	Description  string
	Category     ServiceCategory
	Neighborhood string
}

type ServiceCategory struct {
	Name        string
	Subcategory string
}

func (service Service) ToString() string {
	return fmt.Sprintf(
		"{category: %s,user_name: %s, description: %s, neighborhood: %s}\n",
		service.Category.ToString(),
		service.UserName,
		service.Description,
		service.Neighborhood)
}

func (category ServiceCategory) ToString() string {
	return fmt.Sprintf(
		"{name: %s,subcategory: %s}",
		category.Name,
		category.Subcategory)
}

func (service Service) ToDbFormat() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["user_name"] = service.UserName
	dict["description"] = service.Description
	dict["category"] = service.Category.ToDbFormat()
	if service.Neighborhood != "" {
		dict["neighborhood"] = service.Neighborhood
	}

	return dict
}

func (category ServiceCategory) ToDbFormat() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["name"] = category.Name
	if category.Subcategory != "" {
		dict["subcategory"] = category.Subcategory
	}
	return dict
}

func ServiceIndexes() []mgo.Index {
	indexes := make([]mgo.Index, 4)

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"user_name"},
			Unique: true})

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"category.name"},
			Unique: true})

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"category.subcategory"},
			Unique: true,
			Sparse: true})

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"$text:description"},
			Unique: true})

	return indexes
}
