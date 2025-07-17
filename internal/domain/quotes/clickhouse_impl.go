package quotes

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickhouseQuotesRepositoryImpl struct {
	ClickhouseConn driver.Conn
	Logger         *slog.Logger
}

func NewClickhouseQuotesRepositoryImpl(logger *slog.Logger) *ClickhouseQuotesRepositoryImpl {
	// TODO: HARD dependency, should be injected, but for the purposes of the test, it's fine :)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"192.168.1.7:9000"},
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
		os.Exit(1)
	}

	if err := conn.Ping(context.Background()); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			logger.Error("failed to ping clickhouse",
				slog.Any("code", exception.Code),
				slog.String("message", exception.Message),
				slog.String("stacktrace", exception.StackTrace),
			)
		}
		os.Exit(1)
	}

	return &ClickhouseQuotesRepositoryImpl{
		ClickhouseConn: conn,
		Logger:         logger,
	}
}

func (r *ClickhouseQuotesRepositoryImpl) AddQuotes(ctx context.Context, quotes []Quote) error {
	batch, err := r.ClickhouseConn.PrepareBatch(ctx, "INSERT INTO quotes (id, name, service, deadline, price)")
	if err != nil {
		return errors.Join(errors.New("failed to prepare batch"), err)
	}
	defer func() {
		if err := batch.Close(); err != nil {
			r.Logger.Error("failed to close batch",
				slog.Any("error", err),
			)
		}
	}()

	for _, quote := range quotes {
		if err := batch.Append(
			quote.ID,
			quote.Name,
			quote.Service,
			quote.Deadline,
			quote.Price,
		); err != nil {
			if err := batch.Abort(); err != nil {
				r.Logger.Error("failed to close abort batch",
					slog.Any("error", err),
				)
			}

			return errors.Join(errors.New("failed to append to batch"), err)
		}
	}

	if err := batch.Send(); err != nil {
		r.Logger.Error("failed to close send batch",
			slog.Any("error", err),
		)
		return errors.Join(errors.New("failed to send batch"), err)
	}

	return nil
}
