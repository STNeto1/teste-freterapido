package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
	"github.com/stneto1/teste-freterapido/internal/utils"
)

func (router *Router) QuotesHandler(w http.ResponseWriter, r *http.Request) {
	var body quotes.RequestQuote
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		if _, err := w.Write([]byte(err.Error())); err != nil {
			router.logger.Error(
				"failed to parse body",
				slog.String("error", err.Error()),
			)
		}

	}

	// err => QuoteRequestErrorSetError, QuoteRequestError
	result, err := router.quoteSvc.GetFreteRapidoQuotes(r.Context(), &body)
	if err != nil {
		router.logger.Error("failed to handle request",
			slog.Any("error", err),
			slog.String("err type", reflect.TypeOf(err).String()),
		)

		if _, ok := err.(quotes.QuoteRequestErrorSetError); ok {
			w.WriteHeader(http.StatusBadRequest)

			errResponse := ResponseQuoteErrorSet{
				Errors: err.(quotes.QuoteRequestErrorSetError).Errors,
			}

			if _, err := w.Write(errResponse.ToJSON()); err != nil {
				router.logger.Error(
					"failed to send response",
					slog.String("error", err.Error()),
				)
			}
			return
		}

		if _, ok := err.(quotes.QuoteRequestError); ok {
			w.WriteHeader(http.StatusFailedDependency)

			if _, err := w.Write(ResponseGenericError{
				Message: "External service failed",
			}.ToJSON()); err != nil {
				router.logger.Error(
					"failed to send response",
					slog.String("error", err.Error()),
				)
			}
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(ResponseGenericError{
			Message: "Failed to process request",
		}.ToJSON()); err != nil {
			router.logger.Error(
				"failed to send response",
				slog.String("error", err.Error()),
			)
		}

		return
	}

	response := ResponseQuotes{
		Carrier: utils.Map(result, func(quote quotes.Quote) ResponseQuotesItem {
			return ResponseQuotesItem{
				Name:     quote.Name,
				Service:  quote.Service,
				Deadline: quote.Deadline,
				Price:    quote.Price.InexactFloat64(),
			}
		}),
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		router.logger.Error("failed to marshal response",
			slog.String("error", err.Error()),
		)

		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write(ResponseGenericError{
			Message: "Failed to respond",
		}.ToJSON()); err != nil {
			router.logger.Error(
				"failed to send response",
				slog.String("error", err.Error()),
			)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseBytes); err != nil {
		router.logger.Error("failed to write response",
			slog.String("error", err.Error()),
		)
	}

}
