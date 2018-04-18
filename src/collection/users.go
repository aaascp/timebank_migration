package collection

import "fmt"

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
