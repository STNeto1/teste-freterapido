package analytics

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const baseQuery = `WITH last_quotes AS (SELECT *
                     FROM quotes
                     ORDER BY timestamp DESC
                     %s)
SELECT name       AS carrier,
       count()    AS total_quotes,
       sum(price) AS total_price,
       avg(price) AS avg_price,
       min(price) as min_price,
       max(price) as max_price
FROM last_quotes
GROUP BY carrier
ORDER BY carrier;`

type ClickhouseAnalyticsRepositoryImpl struct {
	ClickhouseConn driver.Conn
	Logger         *slog.Logger
}

func NewClickhouseAnalyticsRepositoryImpl(logger *slog.Logger) *ClickhouseAnalyticsRepositoryImpl {
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

	return &ClickhouseAnalyticsRepositoryImpl{
		ClickhouseConn: conn,
		Logger:         logger,
	}
}

// GetMetrics queries the database for the metrics about the last N requests, if lastQuotes = 0, no limit is applied
func (r *ClickhouseAnalyticsRepositoryImpl) GetMetrics(ctx context.Context, lastQuotes uint64) ([]ServiceMetrics, error) {
	var base string
	params := []any{}
	if lastQuotes > 0 {
		base = fmt.Sprintf(baseQuery, fmt.Sprintf("LIMIT %d", lastQuotes))
		params = append(params, lastQuotes)
	} else {
		base = fmt.Sprintf(baseQuery, "")
	}

	rows, err := r.ClickhouseConn.Query(ctx, base, params...)
	if err != nil {
		return nil, errors.Join(errors.New("failed to query"), err)
	}

	var metrics []ServiceMetrics
	for rows.Next() {
		var metric ServiceMetrics
		if err := rows.ScanStruct(&metric); err != nil {

			return nil, errors.Join(errors.New("failed to scan"), err)
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}
