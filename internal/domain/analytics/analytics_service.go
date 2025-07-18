package analytics

import (
	"context"
	"strconv"

	"github.com/stneto1/teste-freterapido/internal/domain/system"
)

type AnalyticService struct {
	config               *system.AnalyticsServiceConfig
	clickhouseRepository ClickhouseAnalyticsRepository
}

func NewAnalyticService(
	config *system.AnalyticsServiceConfig,
	clickhouseRepository ClickhouseAnalyticsRepository,
) *AnalyticService {
	return &AnalyticService{
		config:               config,
		clickhouseRepository: clickhouseRepository,
	}
}

func (s *AnalyticService) GetAnalytics(ctx context.Context, lastQuotes string) ([]ServiceMetrics, error) {
	if lastQuotes == "" {
		return s.clickhouseRepository.GetMetrics(ctx, 0)
	}

	numericValue, err := strconv.ParseInt(lastQuotes, 10, 64)
	if err != nil {
		return nil, AnalyticsInvalidLastQuotesError{
			Message: "Optional parameter last_quotes must be a number",
		}
	}

	if numericValue < 0 {
		return nil, AnalyticsInvalidLastQuotesError{
			Message: "Optional parameter last_quotes must be greater than 0",
		}
	}

	return s.clickhouseRepository.GetMetrics(ctx, uint64(numericValue))
}
