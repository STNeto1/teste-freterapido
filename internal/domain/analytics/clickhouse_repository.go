package analytics

import "context"

type ClickhouseAnalyticsRepository interface {
	GetMetrics(_ context.Context, LastQuotes uint64) ([]ServiceMetrics, error)
}
