package actions

import (
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"

	"github.com/stellar/go/services/horizon/internal/assets"
	"github.com/stellar/go/xdr"
)

// Validateable allow structs to define their own custom validations.
type Validateable interface {
	Validate() error
}

func InitValidators() {
	govalidator.TagMap["accountID"] = govalidator.Validator(isAccountID)
	govalidator.TagMap["assetType"] = govalidator.Validator(isAssetType)
	govalidator.TagMap["cursor"] = govalidator.Validator(isValidCursor)
	govalidator.CustomTypeTagMap.Set("sellingCode", isValidSellingCode)
	govalidator.CustomTypeTagMap.Set("buyingCode", isValidBuyingCode)

	// govalidator handles embedded structs as fields and tries to find a
	// validator for them, failing if it doesn't find one. This relaxes this
	// setup a bit, so it doesn't try to validate on the embedded struct.
	govalidator.SetFieldsRequiredByDefault(false)
}

func validatorErrorFor(msg string) string {
	switch {
	case strings.HasSuffix(msg, "does not validate as accountID"):
		return "Account ID must start with `G` and contain 56 alphanum characters."
	case strings.HasSuffix(msg, "does not validate as assetType"):
		return "Asset type must be native, credit_alphanum4 or credit_alphanum12."
	case strings.HasSuffix(msg, "does not validate as sellingCode") || strings.HasSuffix(msg, "does not validate as buyingCode"):
		return "Asset code must be 1-12 alphanumeric characters."
	default:
		return msg
	}
}

func isValidCursor(str string) bool {
	// If cursor is a negative value, return false
	cursorInt, err := strconv.Atoi(str)
	if err == nil && cursorInt < 0 {
		return false
	}

	return true
}

func isValidSellingCode(i interface{}, context interface{}) bool {
	switch v := context.(type) {
	case SellingBuyingAssetQueryParams:
		// if asset type is not specified then asset code shouldn't be passed
		if len(v.SellingAssetType) == 0 {
			return false
		}

		t, err := assets.Parse(v.SellingAssetType)
		if err != nil {
			return false
		}

		var validLen int
		switch t {
		case xdr.AssetTypeAssetTypeCreditAlphanum4:
			validLen = len(xdr.AssetAlphaNum4{}.AssetCode)
		case xdr.AssetTypeAssetTypeCreditAlphanum12:
			validLen = len(xdr.AssetAlphaNum12{}.AssetCode)
		}

		if len(v.SellingAssetCode) > validLen {
			return false
		}
	}

	return true
}

func isValidBuyingCode(i interface{}, context interface{}) bool {
	switch v := context.(type) {
	case SellingBuyingAssetQueryParams:
		// if asset type is not specified then asset code shouldn't be passed
		if len(v.BuyingAssetType) == 0 {
			return false
		}

		t, err := assets.Parse(v.BuyingAssetType)
		if err != nil {
			return false
		}

		var validLen int
		switch t {
		case xdr.AssetTypeAssetTypeCreditAlphanum4:
			validLen = len(xdr.AssetAlphaNum4{}.AssetCode)
		case xdr.AssetTypeAssetTypeCreditAlphanum12:
			validLen = len(xdr.AssetAlphaNum12{}.AssetCode)
		}

		if len(v.BuyingAssetCode) > validLen {
			return false
		}
	}

	return true
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
