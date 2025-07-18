package system

import (
	"log/slog"
	"time"
)

type QuotesServiceConfig struct {
	RegisteredNumber  string
	Token             string
	PlatformCode      string
	DispatcherZipCode int64
	Logger            *slog.Logger
	TryQuotesRetries  int
	TryQuotesTimeout  time.Duration
	AddQuotesRetries  int
	AddQuotesTimeout  time.Duration
}

type AnalyticsServiceConfig struct {
	Logger *slog.Logger
}
