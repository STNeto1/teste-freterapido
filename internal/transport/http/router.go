package http

import (
	"log/slog"
	"net/http"

	"github.com/stneto1/teste-freterapido/internal/domain/analytics"
	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
)

type Router struct {
	http.ServeMux
	quoteSvc     *quotes.QuoteService
	analyticsSvc *analytics.AnalyticService
}

func NewRouter(logger *slog.Logger,
	quoteSvc *quotes.QuoteService,
	analyticsSvc *analytics.AnalyticService,
) *Router {

	router := Router{
		ServeMux:     *http.NewServeMux(),
		quoteSvc:     quoteSvc,
		analyticsSvc: analyticsSvc,
	}

	router.HandleFunc("/health", router.HealthHandler)

	return &router
}
