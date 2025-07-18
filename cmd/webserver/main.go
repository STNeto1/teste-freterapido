package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
	"github.com/stneto1/teste-freterapido/internal/domain/analytics"
	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
	"github.com/stneto1/teste-freterapido/internal/domain/system"
	wshttp "github.com/stneto1/teste-freterapido/internal/transport/http"
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

	if kongCtx.Command() != "start" {
		logger.Error("invalid command",
			slog.String("command", kongCtx.Command()),
		)
		os.Exit(1)
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

	router := wshttp.NewRouter(logger,
		quoteSvc,
		analyticsSvc,
	)

	logger.Info("starting webserver",
		slog.String("addr", cliArgs.Start.HTTPAddr),
	)
	if err := http.ListenAndServe(cliArgs.Start.HTTPAddr, router); err != nil {
		logger.Error("server error",
			slog.String("error", err.Error()),
		)
	}
}
