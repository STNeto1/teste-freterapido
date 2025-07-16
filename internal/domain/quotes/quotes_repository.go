package quotes

type ClickhouseQuotesRepository interface {
	CreateQuote(Quote) error
}
