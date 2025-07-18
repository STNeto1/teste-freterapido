package system

import (
	"log/slog"
	"time"
)

type Start struct {
	HTTPAddr          string        `help:"Value to be used as the HTTP addr" default:":8080"`
	RegisteredNumber  string        `help:"Value to be used as the registered number"`
	Token             string        `help:"Value to be used as the token"`
	PlatformCode      string        `help:"Value to be used as the platform code"`
	DispatcherZipCode int64         `help:"Value to be used as the dispatcher zip code"`
	TryQuotesRetries  int           `help:"Value to be used as the number of retries for the try quotes request" default:"3"`
	TryQuotesTimeout  time.Duration `help:"Value to be used as the timeout for the try quotes request" default:"100ms"`
	AddQuotesRetries  int           `help:"Value to be used as the number of retries for the add quotes request" default:"3"`
	AddQuotesTimeout  time.Duration `help:"Value to be used as the timeout for the add quotes request" default:"100ms"`
}

type CLI struct {
	Start Start `cmd:"" help:"Start server"`
}

func (s *Start) ProcessQuoteServiceConfig(logger *slog.Logger) QuotesServiceConfig {
	return QuotesServiceConfig{
		RegisteredNumber:  s.RegisteredNumber,
		Token:             s.Token,
		PlatformCode:      s.PlatformCode,
		DispatcherZipCode: s.DispatcherZipCode,
		Logger:            logger,
		TryQuotesRetries:  s.TryQuotesRetries,
		TryQuotesTimeout:  s.TryQuotesTimeout,
		AddQuotesRetries:  s.AddQuotesRetries,
		AddQuotesTimeout:  s.AddQuotesTimeout,
	}
}
