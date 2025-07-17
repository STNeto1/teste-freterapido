package quotes

import "context"

type ClickhouseQuotesRepository interface {
	AddQuotes(context.Context, []Quote) error
}
