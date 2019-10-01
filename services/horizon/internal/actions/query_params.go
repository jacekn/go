package actions

import (
	"strconv"

	"github.com/stellar/go/services/horizon/internal/db2"
	"github.com/stellar/go/xdr"
)

// PageQueryParams query struct for pagination params
type PageQueryParams struct {
	Cursor string `schema:"cursor" valid:"-"`
	Order  string `schema:"order" valid:"in(asc|desc)~valid values are asc or desc"`
	Limit  string `schema:"limit" valid:"int,range(1|200)~value should be between 1 and 200"`
}

func (q PageQueryParams) getLimit() uint64 {
	var (
		def   = db2.DefaultPageSize
		limit = q.Limit
	)

	if limit == "" {
		return uint64(def)
	}

	asI64, err := strconv.ParseInt(limit, 10, 64)

	if err != nil {
		// should have been valid
		panic(err)
	}

	// TODO: move this validation to validator
	// if asI64 <= 0 {
	// 	err = errors.New("invalid limit: non-positive value provided")
	// } else if asI64 > int64(max) {
	// 	err = errors.Errorf("invalid limit: value provided that is over limit max of %d", max)
	// }

	return uint64(asI64)
}

// PageQuery returns the page query.
func (q PageQueryParams) PageQuery() db2.PageQuery {
	// TODO: add GetCursor
	pageQuery, err := db2.NewPageQuery(q.Cursor, true, q.Order, q.getLimit())

	if err != nil {
		// return pageQuery, problem.MakeInvalidFieldProblem(
		// 	"pagination parameters",
		// 	err,
		// )
		// this should have been validated before
		panic(err)
	}

	return pageQuery
}

// SellingBuyingAssetQueryParams query struct for end-points requiring a selling
// and buying asset
type SellingBuyingAssetQueryParams struct {
	SellingAssetType   string `schema:"selling_asset_type" valid:"assetType~valid types are native credit_alphanum4 or credit_alphanum12"` // TODO using a comma doesn't work in custom message, figure how to make it work
	SellingAssetIssuer string `schema:"selling_asset_issuer" valid:"accountID"`
	SellingAssetCode   string `schema:"selling_asset_code" valid:"sellingCode~code too long"`
	BuyingAssetType    string `schema:"buying_asset_type" valid:"assetType~~valid types are native credit_alphanum4 or credit_alphanum12"`
	BuyingAssetIssuer  string `schema:"buying_asset_issuer" valid:"accountID"`
	BuyingAssetCode    string `schema:"buying_asset_code" valid:"buyingCode~code too long"`
}

// Selling returns an xdr.Asset representing the selling side of the offer.
func (q SellingBuyingAssetQueryParams) Selling() *xdr.Asset {
	if len(q.SellingAssetType) == 0 {
		return nil
	}

	selling, err := BuildAsset(
		q.SellingAssetType,
		q.SellingAssetIssuer,
		q.SellingAssetCode,
	)

	if err != nil {
		panic(err)
	}

	return &selling
}

// Buying returns an *xdr.Asset representing the buying side of the offer.
func (q SellingBuyingAssetQueryParams) Buying() *xdr.Asset {
	if len(q.BuyingAssetType) == 0 {
		return nil
	}

	buying, err := BuildAsset(
		q.BuyingAssetType,
		q.BuyingAssetIssuer,
		q.BuyingAssetCode,
	)

	if err != nil {
		panic(err)
	}

	return &buying
}
