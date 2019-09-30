package actions

import (
	"github.com/asaskevich/govalidator"
)

// Validateable allow structs to define their own custom validations.
type Validateable interface {
	Validate() error
}

func InitValidators() {
	govalidator.TagMap["accountID"] = govalidator.Validator(isAccountID)
}

func isAccountID(str string) bool {
	if _, err := buildAccountID(str); err != nil {
		return false
	}

	return true
}
