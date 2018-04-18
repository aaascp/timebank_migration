package migration

import (
	"strings"
	"time"

	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
)

type UsersList []interface{}
type CreditsList []interface{}

func Users(filename string) (UsersList, CreditsList) {
	const (
		start = 9
		end   = 3961
	)

	file, fileError := xlsx.OpenFile("./resources/btf.xlsx")
	if fileError != nil {
		panic("Error opening file")
	}

	sheet := file.Sheets[2]

	users_list := make(UsersList, 0, end)
	credits_list := make(CreditsList, 0, end)

	rows := sheet.Rows[start:end]
	for _, row := range rows {
		account, accountError := row.Cells[0].FormattedValue()
		name, nameError := row.Cells[1].FormattedValue()
		credits, credits_error := row.Cells[9].Float()

		if account != "" && name != "" && accountError == nil && nameError == nil {
			users_list = append(
				users_list,
				bson.M{
					"account": strings.TrimSpace(account),
					"name":    strings.TrimSpace(name)})
		}

		if credits_error == nil {
			credits_list = append(
				credits_list,
				bson.M{
					"account":   strings.TrimSpace(account),
					"value":     int64(credits * 10),
					"type":      "initial_credit",
					"creted_at": time.Now().Unix()})
		}
	}
	return users_list, credits_list
}
