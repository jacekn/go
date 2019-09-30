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
