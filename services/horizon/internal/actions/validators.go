package actions

import (
	"github.com/asaskevich/govalidator"

	"github.com/stellar/go/services/horizon/internal/assets"
)

// Validateable allow structs to define their own custom validations.
type Validateable interface {
	Validate() error
}

func InitValidators() {
	govalidator.TagMap["accountID"] = govalidator.Validator(isAccountID)
	govalidator.TagMap["assetType"] = govalidator.Validator(isAssetType)
}

func isAssetType(str string) bool {
	if _, err := assets.Parse(str); err != nil {
		return false
	}

	return true
}

func isAccountID(str string) bool {
	if _, err := buildAccountID(str); err != nil {
		return false
	}

	return true
}
