package collection

import "fmt"

type Service struct {
	UserName     string
	Description  string
	Category     ServiceCategory
	Neighborhood string
}

func (service Service) ToString() string {
	return fmt.Sprintf(
		"{category: %s,user_name: %s, description: %s, neighborhood: %s}\n",
		service.Category.ToString(),
		service.UserName,
		service.Description,
		service.Neighborhood)
}
