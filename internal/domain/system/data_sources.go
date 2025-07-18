package system

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func CreateClickhouseDatasource(logger *slog.Logger, clickhouseAddr string) (driver.Conn, error) {
	// TODO: HARD dependency, should be injected, but for the purposes of the test, it's fine :)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: strings.Split(clickhouseAddr, ","),
		Auth: clickhouse.Auth{
			Database: "freterapido",
			Username: "default",
			Password: "admin",
		},
		Debug: true,
	})

	if err != nil {
		logger.Error("failed to connect to clickhouse",
			slog.Any("error", err),
		)
		return nil, errors.Join(errors.New("failed to connect to clickhouse"), err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			logger.Error("failed to ping clickhouse",
				slog.Any("code", exception.Code),
				slog.String("message", exception.Message),
				slog.String("stacktrace", exception.StackTrace),
			)
		}
		return nil, errors.Join(errors.New("failed to ping clickhouse"), err)
	}

	return conn, nil
}
