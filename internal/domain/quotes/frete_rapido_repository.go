package quotes

import "context"

// FreteRapidoQuotesRepository - TryQuote sends a request to the Frete Rápido API and returns the response.
// DOES NOT RETRY, you must handle retries yourself.
type FreteRapidoQuotesRepository interface {
	TryQuote(context.Context, FreteRapidoRequestQuote) (FreteRapidoResponseQuote, error)
}
