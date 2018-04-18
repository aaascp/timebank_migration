package migration

import (
	"fmt"
	"strings"
	"timebank/src/collection"

	"github.com/tealeg/xlsx"
)

type Category struct {
	name          string
	subcategories []string
}

type ServicesCollection []collection.Collection
type CategoriesCollection []collection.Collection

var services []collection.Collection
var categories []collection.Collection

var categoryObject collection.ServiceCategory
var lastCategory collection.Category

func Services(filename string) (ServicesCollection, CategoriesCollection) {
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
		isMergedCell := row.Cells[0].HMerge > 1
		isBoldSecond := row.Cells[1].GetStyle().Font.Bold

		first, _ := row.Cells[0].FormattedValue()
		second, _ := row.Cells[1].FormattedValue()
		third, _ := row.Cells[2].FormattedValue()

		if isCategory(isBoldSecond, isMergedCell, first, second, third) {
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
	if lastCategory.Name != "" {
		categories = append(categories, lastCategory)
	}

	categoryName := strings.TrimSpace(category[1])

	categoryObject.Name = categoryName
	categoryObject.Subcategory = ""

	lastCategory.Name = categoryName
	lastCategory.Subcategories = make([]string, 0, 5)
}

func addSubcategory(category string) {
	subcategory := strings.TrimSpace(category)
	categoryObject.Subcategory = subcategory
	lastCategory.Subcategories = append(lastCategory.Subcategories, subcategory)
}

func addService(name, description, neighborhood string) {
	if name == "" || description == "" {
		return
	}

	var service collection.Service
	service.UserName = strings.TrimSpace(name)
	service.Description = strings.TrimSpace(description)
	service.Category = categoryObject
	service.Neighborhood = strings.TrimSpace(neighborhood)

	services = append(services, service)
}

func isCategory(isBoldSecond, isMerged bool, first, second, third string) bool {
	/*
			** There is a row with description only in category 4 - CASA / HORTA / JARDIM
			** 'Já participei com projetos de horta comunitária, agricultura sustentável e sistemas agroflorestais'
			** 'isBold' is a workaround
		  **
			** There is a row with UserName only in category 4 - CASA / HORTA / JARDIM
			** 'isMerged' is a workarounds
	*/
	return ((first != "" && second == "" && third == "" && isMerged) ||
		(first == "" && second != "" && third == "" && isBoldSecond))
}

func categoryName(first, second string) []string {
	if first != "" {
		return strings.Split(first, "-")
	}
	return strings.Split(second, "-")
}
