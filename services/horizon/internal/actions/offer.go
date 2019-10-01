package actions

import (
	"context"
	"net/http"

	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/services/horizon/internal/resourceadapter"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/render/hal"
)

// OffersQuery query struct for offers end-point
type OffersQuery struct {
	PageQueryParams
	SellingBuyingAssetQueryParams
	Seller string `schema:"seller" valid:"accountID"`
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
		PageQuery: qp.PageQuery(r),
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
