package actions

import (
	"context"
	"net/http"
	"strconv"

	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/db2"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/services/horizon/internal/resourceadapter"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/render/hal"
	"github.com/stellar/go/xdr"
)

// QueryParams query struct for pagination params
// TODO: move a shared package - maybe even db2
type QueryParams struct {
	Cursor string `schema:"cursor"`
	Order  string `schema:"order"`
	Limit  string `schema:"limit"` // todo validate uint64
}

func (q QueryParams) getLimit() uint64 {
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
func (q QueryParams) PageQuery() db2.PageQuery {
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

// OffersQuery query struct for offers end-point
type OffersQuery struct {
	QueryParams
	Seller              string `schema:"seller"`
	SellingAssetType    string `schema:"selling_asset_type"`
	SellingAsssetIssuer string `schema:"selling_asset_issuer"`
	SellingAsssetCode   string `schema:"selling_asset_code"`
	BuyingAssetType     string `schema:"buying_asset_type"`
	BuyingAsssetIssuer  string `schema:"buying_asset_issuer"`
	BuyingAsssetCode    string `schema:"buying_asset_code"`
}

// Selling an xdr.Asset representing the selling side of the offer.
func (q OffersQuery) Selling() *xdr.Asset {
	if len(q.SellingAssetType) == 0 {
		return nil
	}

	selling, err := BuildAsset(q.SellingAssetType, q.SellingAsssetIssuer, q.SellingAsssetCode)

	if err != nil {
		panic(err)
	}

	return &selling
}

// Buying an xdr.Asset representing the buying side of the offer.
func (q OffersQuery) Buying() *xdr.Asset {
	if len(q.SellingAssetType) == 0 {
		return nil
	}

	buying, err := BuildAsset(q.BuyingAssetType, q.BuyingAsssetIssuer, q.BuyingAsssetCode)

	if err != nil {
		panic(err)
	}

	return &buying
}

// GetOffersHandler is the action handler for the /offers endpoint
type GetOffersHandler struct {
	HistoryQ *history.Q
}

// GetResourcePage returns a page of offers.
func (handler GetOffersHandler) GetResourcePage(r *http.Request) ([]hal.Pageable, error) {
	ctx := r.Context()
	qp := OffersQuery{}
	err := GetParams(&qp, r)

	if err != nil {
		return nil, err
	}

	query := history.OffersQuery{
		PageQuery: qp.PageQuery(),
		SellerID:  qp.Seller,
		Selling:   qp.Selling(),
		Buying:    qp.Buying(),
	}

	offers, err := getOffersPage(ctx, handler.HistoryQ, query)
	if err != nil {
		return nil, err
	}

	return offers, nil
}

// GetAccountOffersHandler is the action handler for the
// `/accounts/{account_id}/offers` endpoint when using experimental ingestion.
type GetAccountOffersHandler struct {
	HistoryQ *history.Q
}

func (handler GetAccountOffersHandler) parseOffersQuery(r *http.Request) (history.OffersQuery, error) {
	pq, err := GetPageQuery(r)
	if err != nil {
		return history.OffersQuery{}, err
	}

	seller, err := GetString(r, "account_id")
	if err != nil {
		return history.OffersQuery{}, err
	}

	query := history.OffersQuery{
		PageQuery: pq,
		SellerID:  seller,
	}

	return query, nil
}

// GetResourcePage returns a page of offers for a given account.
func (handler GetAccountOffersHandler) GetResourcePage(r *http.Request) ([]hal.Pageable, error) {
	ctx := r.Context()
	query, err := handler.parseOffersQuery(r)
	if err != nil {
		return nil, err
	}

	offers, err := getOffersPage(ctx, handler.HistoryQ, query)
	if err != nil {
		return nil, err
	}

	return offers, nil
}

func getOffersPage(ctx context.Context, historyQ *history.Q, query history.OffersQuery) ([]hal.Pageable, error) {
	records, err := historyQ.GetOffers(query)
	if err != nil {
		return nil, err
	}

	ledgerCache := history.LedgerCache{}
	for _, record := range records {
		ledgerCache.Queue(int32(record.LastModifiedLedger))
	}

	if err := ledgerCache.Load(historyQ); err != nil {
		return nil, errors.Wrap(err, "failed to load ledger batch")
	}

	var offers []hal.Pageable
	for _, record := range records {
		var offerResponse horizon.Offer

		ledger, found := ledgerCache.Records[int32(record.LastModifiedLedger)]
		ledgerPtr := &ledger
		if !found {
			ledgerPtr = nil
		}

		resourceadapter.PopulateHistoryOffer(ctx, &offerResponse, record, ledgerPtr)
		offers = append(offers, offerResponse)
	}

	return offers, nil
}
