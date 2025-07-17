package system

import "log/slog"

type QuotesServiceConfig struct {
	RegisteredNumber  string
	Token             string
	PlatformCode      string
	DispatcherZipCode int64
	Logger            *slog.Logger
}

type AnalyticsServiceConfig struct {
	Logger *slog.Logger
}
