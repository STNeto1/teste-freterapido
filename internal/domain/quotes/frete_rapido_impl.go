package quotes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/stneto1/teste-freterapido/internal/domain/system"
)

type FreteRapidoQuotesRepositoryImpl struct {
	quoteURLEndpoint string
	Logger           *slog.Logger
}

func NewFreteRapidoQuotesRepositoryImpl(logger *slog.Logger) *FreteRapidoQuotesRepositoryImpl {
	// TODO: make it configurable in the future
	// enabling to do a more "e2e" test
	return &FreteRapidoQuotesRepositoryImpl{
		quoteURLEndpoint: "https://sp.freterapido.com/api/v3/quote/simulate",
		Logger:           logger,
	}
}

// TryQuote decode sends a request to the Frete Rápido API and returns the response. DOES NOT RETRY, you must handle it.
func (r *FreteRapidoQuotesRepositoryImpl) TryQuote(ctx context.Context, requestQuote FreteRapidoRequestQuote) (FreteRapidoResponseQuote, error) {
	payload, err := json.Marshal(requestQuote)
	if err != nil {
		return FreteRapidoResponseQuote{}, system.SystemError{
			Message: fmt.Errorf("failed to serialize request\n%s", err.Error()).Error(),
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		r.quoteURLEndpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		return FreteRapidoResponseQuote{}, system.SystemError{
			Message: fmt.Errorf("failed to create request\n%s", err.Error()).Error(),
		}
	}

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FreteRapidoResponseQuote{}, system.SystemError{
			Message:     fmt.Errorf("failed to make request\n%s", err.Error()).Error(),
			ShouldRetry: true,
		}
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			r.Logger.Error("failed to close response body",
				slog.Any("error", err),
			)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return FreteRapidoResponseQuote{}, system.SystemError{
			Message: fmt.Errorf("failed read response\n%s", err.Error()).Error(),
		}
	}

	var responseBody FreteRapidoResponseQuote
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return FreteRapidoResponseQuote{}, system.SystemError{
			Message: fmt.Errorf("failed decode response\n%s", err.Error()).Error(),
		}
	}

	return responseBody, nil
}
