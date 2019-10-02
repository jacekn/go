package actions

import (
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"

	"github.com/stellar/go/services/horizon/internal/db2"
	"github.com/stellar/go/services/horizon/internal/ledger"
	"github.com/stellar/go/services/horizon/internal/toid"
	"github.com/stellar/go/support/render/problem"
)

func TestPageQueryParamOrderValidation(t *testing.T) {
	InitValidators()
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
	InitValidators()
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

func TestPageQueryParamCursorValidation(t *testing.T) {
	InitValidators()
	for _, testCase := range []struct {
		value string
		valid bool
	}{
		{
			"1",
			true,
		},
		{
			"0",
			true,
		},
		{
			"string",
			true,
		},
		{
			"-1",
			false,
		},
	} {
		t.Run(testCase.value, func(t *testing.T) {
			tt := assert.New(t)

			q := PageQueryParams{
				Cursor: testCase.value,
			}

			result, err := govalidator.ValidateStruct(q)
			if testCase.valid {
				tt.NoError(err)
				tt.True(result)
			} else {
				tt.Equal("the value should not be a negative number", err.Error())
			}
		})
	}
}

func TestPageQueryParamGetLimit(t *testing.T) {
	InitValidators()
	tt := assert.New(t)
	q := PageQueryParams{
		Limit: "",
	}

	tt.Equal(uint64(db2.DefaultPageSize), q.getLimit())

	q = PageQueryParams{
		Limit: "20",
	}

	tt.Equal(uint64(20), q.getLimit())
}

func TestPageQueryParamGetCursor(t *testing.T) {
	InitValidators()
	tt := assert.New(t)
	q := PageQueryParams{}
	r := makeAction("/", nil).R

	err := GetParams(&q, r)
	tt.NoError(err)
	tt.Equal("", q.getCursor(r))

	r = makeAction("/?cursor=12345", nil).R

	err = GetParams(&q, r)
	tt.NoError(err)
	tt.Equal("12345", q.getCursor(r))

	r = makeAction("/?cursor=-1", nil).R

	err = GetParams(&q, r)
	tt.Error(err)

	r = makeAction("/?cursor=now", nil).R
	err = GetParams(&q, r)
	tt.NoError(err)
	expected := toid.AfterLedger(ledger.CurrentState().HistoryLatest).String()
	tt.Equal(expected, q.getCursor(r))

	//Last-Event-ID overrides cursor
	r = makeAction("/?cursor=now", nil).R
	r.Header.Set("Last-Event-ID", "from_header")
	err = GetParams(&q, r)
	tt.NoError(err)
	tt.Equal("from_header", q.getCursor(r))
}

func TestSellingBuyingAssetQueryParams(t *testing.T) {
	InitValidators()
	tt := assert.New(t)

	urlParams := map[string]string{
		"selling_asset_type": "invalid",
	}

	r := makeAction("/", urlParams).R
	qp := SellingBuyingAssetQueryParams{}
	err := GetParams(&qp, r)

	if tt.IsType(&problem.P{}, err) {
		p := err.(*problem.P)
		tt.Equal("bad_request", p.Type)
		tt.Equal("selling_asset_type", p.Extras["invalid_field"])
		tt.Equal(
			"Asset type must be native, credit_alphanum4 or credit_alphanum12.",
			p.Extras["reason"],
		)
	}

	urlParams = map[string]string{
		"selling_asset_type": "credit_alphanum4",
		"selling_asset_code": "invalid",
	}

	r = makeAction("/", urlParams).R
	qp = SellingBuyingAssetQueryParams{}
	err = GetParams(&qp, r)

	if tt.IsType(&problem.P{}, err) {
		p := err.(*problem.P)
		tt.Equal("bad_request", p.Type)
		tt.Equal("selling_asset_code", p.Extras["invalid_field"])
		tt.Equal(
			"Asset code must be 1-12 alphanumeric characters.",
			p.Extras["reason"],
		)
	}
}
