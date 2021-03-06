package collection

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
)

type Transaction struct {
	Account   string
	Value     int64
	Type      string
	CreatedAt int64
}

func (transaction Transaction) ToString() string {
	return fmt.Sprintf(
		"{credit_account: %s, value: %d, type: %s, created_at: %s}\n",
		transaction.Account,
		transaction.Value,
		transaction.Type,
		formatCreatedAt(transaction.CreatedAt))
}

func (transaction Transaction) ToDbFormat() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["credit_account"] = transaction.Account
	dict["value"] = transaction.Value
	dict["type"] = transaction.Type
	dict["created_at"] = transaction.CreatedAt
	return dict
}

func TransactionIndexes() []mgo.Index {
	indexes := make([]mgo.Index, 2)

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"debit_account"},
			Unique: true})

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"credit_account"},
			Unique: true})

	return indexes
}

func formatCreatedAt(createdAt int64) string {
	createdAtTime := time.Unix(createdAt, 0)

	return fmt.Sprintf(
		"%02d:%02d:%02d %02d/%02d/%d",
		createdAtTime.Hour(),
		createdAtTime.Minute(),
		createdAtTime.Second(),
		createdAtTime.Day(),
		createdAtTime.Month(),
		createdAtTime.Year())
}
