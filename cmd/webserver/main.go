package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/shopspring/decimal"
	"github.com/stneto1/teste-freterapido/internal/domain/analytics"
	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
	"github.com/stneto1/teste-freterapido/internal/domain/system"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var cliArgs system.CLI

	kongCtx := kong.Parse(&cliArgs,
		kong.Name("freterapido-webserver"),
		kong.Description("Webserver for frete rapido"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)
	switch kongCtx.Command() {
	case "start":
		logger.Info("starting webserver",
			slog.Any("args", cliArgs),
		)
	default:
		panic(kongCtx.Command())
	}

	quoteCfg := cliArgs.Start.ProcessQuoteServiceConfig(logger)
	analyticsCfg := system.AnalyticsServiceConfig{
		Logger: logger,
	}

	fastFreteSource := quotes.NewFreteRapidoQuotesRepositoryImpl(quoteCfg.Logger)
	clickhouseQuotesSource := quotes.NewClickhouseQuotesRepositoryImpl(quoteCfg.Logger)
	clickhouseAnalyticsSource := analytics.NewClickhouseAnalyticsRepositoryImpl(analyticsCfg.Logger)

	quoteSvc := quotes.NewQuoteService(
		&quoteCfg,
		fastFreteSource,
		clickhouseQuotesSource,
	)
	analyticsSvc := analytics.NewAnalyticService(
		&analyticsCfg,
		clickhouseAnalyticsSource,
	)

	result, err := quoteSvc.GetFreteRapidoQuotes(context.Background(), &quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "29161376",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category:      7,
				Amount:        1,
				UnitaryWeight: 4,
				Price:         decimal.NewFromFloat(349),
				Sku:           "abs-teste-123",
				Height:        0.2,
				Width:         0.2,
				Length:        0.2,
			},
		},
	})

	if err != nil {
		quoteCfg.Logger.Error("failed to get quotes",
			slog.String("error", err.Error()),
		)
		return
	}

	quoteCfg.Logger.Info("got quotes",
		slog.Any("quotes", result),
	)

	// MAKE SURE THE GOROUTINE IS RUNNING
	time.Sleep(time.Second * 2)

	metrics, err := analyticsSvc.GetAnalytics(context.Background(), "10")
	if err != nil {
		analyticsCfg.Logger.Error("failed to get metrics",
			slog.String("error", err.Error()),
		)
		return
	}

	analyticsCfg.Logger.Info("got metrics",
		slog.Any("metrics", metrics),
	)
}
