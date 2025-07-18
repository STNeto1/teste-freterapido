package http

import (
	"log/slog"
	"net/http"

	"github.com/stneto1/teste-freterapido/internal/domain/analytics"
	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
)

type Router struct {
	http.ServeMux
	logger       *slog.Logger
	quoteSvc     *quotes.QuoteService
	analyticsSvc *analytics.AnalyticService
}

func NewRouter(
	logger *slog.Logger,
	quoteSvc *quotes.QuoteService,
	analyticsSvc *analytics.AnalyticService,
) *Router {

	router := Router{
		ServeMux:     *http.NewServeMux(),
		logger:       logger,
		quoteSvc:     quoteSvc,
		analyticsSvc: analyticsSvc,
	}

	router.HandleFunc("GET /health", router.HealthHandler)
	router.HandleFunc("GET /metrics", router.MetricsHandler)
	router.HandleFunc("POST /quotes", router.QuotesHandler)

	return &router
}
