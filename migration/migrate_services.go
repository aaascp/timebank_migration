package migration

import (
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
)

type Category struct {
	name          string
	subcategories []string
}

type CategoryObject struct {
	name        string
	subcategory string
}

type ServicesList []interface{}
type CategoriesList []interface{}

var services ServicesList
var categories CategoriesList

var categoryObject CategoryObject
var lastCategory Category

var isFirstCategory bool

func Services(filename string) (ServicesList, CategoriesList) {
	const (
		start = 3
		end   = 6990
	)

	file, fileError := xlsx.OpenFile(filename)
	if fileError != nil {
		panic(fmt.Sprintf("Error opening file: %s\n", fileError))
	}

	sheet := file.Sheets[1]

	rows := sheet.Rows[start:end]
	for _, row := range rows {
		first, _ := row.Cells[0].FormattedValue()
		second, _ := row.Cells[1].FormattedValue()
		third, _ := row.Cells[2].FormattedValue()

		if isCategory(first, second, third) {
			createCategory(first, second)
		} else {
			addService(first, second, third)
		}
	}

	categories = append(categories, lastCategory)
	return services, categories
}

func createCategory(first, second string) {
	category := categoryName(first, second)
	if len(category) == 2 {
		createTopCategory(category)
	} else {
		addSubcategory(category[0])
	}
}

func createTopCategory(category []string) {
	if !isFirstCategory {
		categories = append(
			categories,
			bson.M{
				"name":          lastCategory.name,
				"subcategories": lastCategory.subcategories})
	} else {
		isFirstCategory = false
	}

	categoryName := strings.TrimSpace(category[1])

	categoryObject.name = categoryName

	lastCategory.name = categoryName
	lastCategory.subcategories = make([]string, 0, 5)
}

func addSubcategory(category string) {
	subcategory := strings.TrimSpace(category)
	categoryObject.subcategory = subcategory
	lastCategory.subcategories = append(lastCategory.subcategories, subcategory)
}

func addService(name, description, neighborhood string) {
	if name == "" {
		return
	}

	category := bson.M{"name": categoryObject.name}

	if categoryObject.subcategory != "" {
		category["subcategory"] = categoryObject.subcategory
	}

	service := bson.M{
		"user_name":   strings.TrimSpace(name),
		"description": strings.TrimSpace(description),
		"category":    category}

	if neighborhood != "" {
		service["neighborhood"] = strings.TrimSpace(neighborhood)
	}

	services = append(services, service)
}

func isCategory(first, second, third string) bool {
	return ((first != "" && second == "" && third == "") ||
		(first == "" && second != "" && third == ""))
}

func categoryName(first, second string) []string {
	if first != "" {
		return strings.Split(first, "-")
	}
	return strings.Split(second, "-")
}
