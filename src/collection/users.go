package collection

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

type User struct {
	Account string
	Name    string
}

func (user User) ToString() string {
	return fmt.Sprintf(
		"{account: %s, name: %s}\n",
		user.Account,
		user.Name)
}

func (user User) ToDbFormat() map[string]interface{} {
	dict := make(map[string]interface{})
	dict["account"] = user.Account
	dict["name"] = user.Name
	return dict
}

func UserIndexes() []mgo.Index {
	indexes := make([]mgo.Index, 3)

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"name"},
			Unique: true})

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"facebook_id"},
			Unique: true,
			Sparse: true})

	indexes = append(indexes,
		mgo.Index{
			Key:    []string{"account"},
			Unique: true})

	return indexes
}
