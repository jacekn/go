package actions

import (
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

func TestPageQueryParamOrderValidation(t *testing.T) {
	for _, testCase := range []struct {
		desc  string
		order string
		valid bool
	}{
		{
			"empty string",
			"",
			true,
		},
		{
			"desc",
			"desc",
			true,
		},
		{
			"asc",
			"asc",
			true,
		},
		{
			"invalid order",
			"foo",
			false,
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			tt := assert.New(t)

			q := PageQueryParams{
				Order: testCase.order,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.valid {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal("valid values are asc or desc", err.Error())
			}
		})
	}
}

func TestPageQueryParamLimitValidation(t *testing.T) {
	for _, testCase := range []struct {
		desc        string
		limit       string
		expectedMsg string
	}{
		{
			"positive value",
			"1",
			"",
		},
		{
			"non-positive",
			"-1",
			"non-positive value provided",
		},
		{
			"zero",
			"0",
			"non-positive value provided",
		},
		{
			"non-numerical value",
			"foo",
			"Limit: foo does not validate as int",
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			tt := assert.New(t)

			q := PageQueryParams{
				Limit: testCase.limit,
			}

			result, err := govalidator.ValidateStruct(q)
			if len(testCase.expectedMsg) == 0 {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal(testCase.expectedMsg, err.Error())
			}
		})
	}
}
