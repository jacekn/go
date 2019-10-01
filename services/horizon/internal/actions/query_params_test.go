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
		desc  string
		limit string
		valid bool
	}{
		{
			"positive value",
			"1",
			true,
		},
		{
			"positive value",
			"200",
			true,
		},
		{
			"positive value",
			"201",
			false,
		},
		{
			"non-positive",
			"-1",
			false,
		},
		{
			"zero",
			"0",
			false,
		},
	} {
		t.Run(testCase.desc, func(t *testing.T) {
			tt := assert.New(t)

			q := PageQueryParams{
				Limit: testCase.limit,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.valid {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal("value should be between 1 and 200", err.Error())
			}
		})
	}
}
