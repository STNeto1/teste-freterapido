package analytics

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/stneto1/teste-freterapido/internal/domain/system"
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

func NewClickhouseAnalyticsRepositoryImpl(logger *slog.Logger, clickhouseAddr string) *ClickhouseAnalyticsRepositoryImpl {
	conn, err := system.CreateClickhouseDatasource(logger, clickhouseAddr)
	if err != nil {
		logger.Error("failed to create clickhouse datasource",
			slog.Any("error", err),
		)
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

	metrics := make([]ServiceMetrics, 0)

	for rows.Next() {
		var metric ServiceMetrics
		if err := rows.ScanStruct(&metric); err != nil {

			return nil, errors.Join(errors.New("failed to scan"), err)
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}
