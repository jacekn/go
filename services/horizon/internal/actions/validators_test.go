package actions

import (
	"os"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	InitValidators()
	code := m.Run()
	os.Exit(code)
}

func TestBuyingCodeValidator(t *testing.T) {
	for _, testCase := range []struct {
		assetType string
		assetCode string
		valid     bool
	}{
		{
			"native",
			"",
			true,
		},
		{
			"credit_alphanum4",
			"USD",
			true,
		},
		{
			"credit_alphanum12",
			"OHLOOONG",
			true,
		},
		{
			"credit_alphanum4",
			"USDXD",
			false,
		},
		{
			"credit_alphanum12",
			"OHLOOOOOOOOOONG",
			false,
		},
	} {
		t.Run(testCase.assetType, func(t *testing.T) {
			tt := assert.New(t)

			q := SellingBuyingAssetQueryParams{
				BuyingAssetType: testCase.assetType,
				BuyingAssetCode: testCase.assetCode,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.valid {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal("code too long", err.Error())
			}
		})
	}
}

func TestSellingCodeValidator(t *testing.T) {
	for _, testCase := range []struct {
		assetType string
		assetCode string
		valid     bool
	}{
		{
			"native",
			"",
			true,
		},
		{
			"credit_alphanum4",
			"USD",
			true,
		},
		{
			"credit_alphanum12",
			"OHLOOONG",
			true,
		},
		{
			"credit_alphanum4",
			"USDXD",
			false,
		},
		{
			"credit_alphanum12",
			"OHLOOOOOOOOOONG",
			false,
		},
	} {
		t.Run(testCase.assetType, func(t *testing.T) {
			tt := assert.New(t)

			q := SellingBuyingAssetQueryParams{
				SellingAssetType: testCase.assetType,
				SellingAssetCode: testCase.assetCode,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.valid {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal("code too long", err.Error())
			}
		})
	}

}

func TestAssetTypeValidator(t *testing.T) {
	type Query struct {
		AssetType string `valid:"assetType~invalid asset type"`
	}

	for _, testCase := range []struct {
		assetType string
		valid     bool
	}{
		{
			"native",
			true,
		},
		{
			"credit_alphanum4",
			true,
		},
		{
			"credit_alphanum12",
			true,
		},
		{
			"",
			true,
		},
		{
			"stellar_asset_type",
			false,
		},
	} {
		t.Run(testCase.assetType, func(t *testing.T) {
			tt := assert.New(t)

			q := Query{
				AssetType: testCase.assetType,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.valid {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal("invalid asset type", err.Error())
			}
		})
	}
}

func TestAccountIDValidator(t *testing.T) {
	type Query struct {
		Account string `valid:"accountID~invalid address"`
	}

	for _, testCase := range []struct {
		name          string
		value         string
		expectedError string
	}{
		{
			"invalid stellar address",
			"FON4WOTCFSASG3J6SGLLQZURDDUVNBQANAHEQJ3PBNDZ74X63UZWQPZW",
			"invalid address",
		},
		{
			"valid stellar address",
			"GAN4WOTCFSASG3J6SGLLQZURDDUVNBQANAHEQJ3PBNDZ74X63UZWQPZW",
			"",
		},
		{
			"empty stellar address should not be validated",
			"",
			"",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			tt := assert.New(t)

			q := Query{
				Account: testCase.value,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.expectedError == "" {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal(testCase.expectedError, err.Error())
			}
		})
	}
}
