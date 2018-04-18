package migration

import (
	"strings"
	"time"
	"timebank/src/collection"

	"github.com/tealeg/xlsx"
)

type UsersList []collection.Collection
type CreditsList []collection.Collection

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

	usersList := make(UsersList, 0, end)
	creditsList := make(CreditsList, 0, end)

	rows := sheet.Rows[start:end]
	for _, row := range rows {
		account, accountError := row.Cells[0].FormattedValue()
		name, nameError := row.Cells[1].FormattedValue()
		credits, creditsError := row.Cells[9].Float()

		if account != "" && name != "" && accountError == nil && nameError == nil {
			usersList = append(
				usersList,
				collection.User{
					strings.TrimSpace(account),
					strings.TrimSpace(name)})
		}

		if creditsError == nil {
			creditsList = append(
				creditsList,
				collection.Transaction{
					strings.TrimSpace(account),
					int64(credits * 10),
					"initial_credit",
					time.Now().Unix()})
		}
	}
	return usersList, creditsList
}
