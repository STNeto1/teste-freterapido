package quotes

import "context"

type ClickhouseQuotesRepositoryImpl struct {
}

func NewClickhouseQuotesRepositoryImpl() *ClickhouseQuotesRepositoryImpl {
	return &ClickhouseQuotesRepositoryImpl{}
}

func (r *ClickhouseQuotesRepositoryImpl) AddQuotes(_ctx context.Context, _quotes []Quote) error {
	// NOOP func
	return nil
}
