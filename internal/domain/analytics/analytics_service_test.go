package analytics_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stneto1/teste-freterapido/internal/domain/analytics"
	"github.com/stneto1/teste-freterapido/internal/domain/system"
	"github.com/stneto1/teste-freterapido/mocks/analyticsmocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAnalyticsService_GetAnalytics_EmptyString(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clickhouseMock := analyticsmocks.NewMockClickhouseAnalyticsRepository(ctrl)
	clickhouseMock.
		EXPECT().
		GetMetrics(gomock.Any(), uint64(0)).
		DoAndReturn(func(_ context.Context, _ uint64) ([]analytics.ServiceMetrics, error) {
			return []analytics.ServiceMetrics{}, nil
		})

	analyticsCfg := &system.AnalyticsServiceConfig{
		Logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}
	svc := analytics.NewAnalyticService(analyticsCfg, clickhouseMock)

	result, err := svc.GetAnalytics(t.Context(), "")
	assert.Equal(t, []analytics.ServiceMetrics{}, result)
	assert.NoError(t, err)
}

func TestAnalyticsService_GetAnalytics_InvalidNumericString(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clickhouseMock := analyticsmocks.NewMockClickhouseAnalyticsRepository(ctrl)
	clickhouseMock.
		EXPECT().
		GetMetrics(gomock.Any(), uint64(0)).
		DoAndReturn(func(_ context.Context, _ uint64) ([]analytics.ServiceMetrics, error) {
			return []analytics.ServiceMetrics{}, nil
		}).
		Times(0)

	analyticsCfg := &system.AnalyticsServiceConfig{
		Logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}
	svc := analytics.NewAnalyticService(analyticsCfg, clickhouseMock)

	result, err := svc.GetAnalytics(t.Context(), "SUT")
	assert.Nil(t, result)
	assert.ErrorIs(t, err, analytics.AnalyticsInvalidLastQuotesError{
		Message: "Optional parameter last_quotes must be a number",
	})
}

func TestAnalyticsService_GetAnalytics_NegativeNumericString(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clickhouseMock := analyticsmocks.NewMockClickhouseAnalyticsRepository(ctrl)
	clickhouseMock.
		EXPECT().
		GetMetrics(gomock.Any(), uint64(0)).
		DoAndReturn(func(_ context.Context, _ uint64) ([]analytics.ServiceMetrics, error) {
			return []analytics.ServiceMetrics{}, nil
		}).
		Times(0)

	analyticsCfg := &system.AnalyticsServiceConfig{
		Logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}
	svc := analytics.NewAnalyticService(analyticsCfg, clickhouseMock)

	result, err := svc.GetAnalytics(t.Context(), "-1")
	assert.Nil(t, result)
	assert.ErrorIs(t, err, analytics.AnalyticsInvalidLastQuotesError{
		Message: "Optional parameter last_quotes must be greater than 0",
	})
}

func TestAnalyticsService_GetAnalytics_ValidNumericString(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clickhouseMock := analyticsmocks.NewMockClickhouseAnalyticsRepository(ctrl)
	clickhouseMock.
		EXPECT().
		GetMetrics(gomock.Any(), uint64(10)).
		DoAndReturn(func(_ context.Context, _ uint64) ([]analytics.ServiceMetrics, error) {
			return []analytics.ServiceMetrics{}, nil
		})

	analyticsCfg := &system.AnalyticsServiceConfig{
		Logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}
	svc := analytics.NewAnalyticService(analyticsCfg, clickhouseMock)

	result, err := svc.GetAnalytics(t.Context(), "10")
	assert.Equal(t, []analytics.ServiceMetrics{}, result)
	assert.NoError(t, err)
}
