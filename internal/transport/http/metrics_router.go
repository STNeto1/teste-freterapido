package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/stneto1/teste-freterapido/internal/domain/analytics"
)

func (router *Router) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	lastQuotes := r.URL.Query().Get("last_quotes")

	result, err := router.analyticsSvc.GetAnalytics(r.Context(), lastQuotes)
	if err != nil {
		if errors.Is(err, analytics.AnalyticsInvalidLastQuotesError{}) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if _, err := w.Write([]byte(err.Error())); err != nil {
			router.logger.Error(
				"failed to send response",
				slog.String("error", err.Error()),
			)
		}

		return
	}

	responseBytes, err := json.Marshal(result)
	if err != nil {
		router.logger.Error("failed to marshal response",
			slog.String("error", err.Error()),
		)

		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Failed to encode response")); err != nil {
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
