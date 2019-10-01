package actions

import (
	"strconv"

	"github.com/asaskevich/govalidator"

	"github.com/stellar/go/services/horizon/internal/assets"
	"github.com/stellar/go/xdr"
)

// Validateable allow structs to define their own custom validations.
type Validateable interface {
	Validate() error
}

func InitValidators() {
	govalidator.TagMap["isPositive"] = govalidator.Validator(isPositive)
	govalidator.TagMap["accountID"] = govalidator.Validator(isAccountID)
	govalidator.TagMap["assetType"] = govalidator.Validator(isAssetType)
	govalidator.CustomTypeTagMap.Set("sellingCode", isValidSellingCode)
	govalidator.CustomTypeTagMap.Set("buyingCode", isValidBuyingCode)
}

func isPositive(str string) bool {
	asI64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return false
	}

	return asI64 > 0
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
