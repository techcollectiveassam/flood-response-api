package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
)

type queryContext struct {
	start time.Time
	sql   string
}

// QueryTracer logs every SQL query executed through the pool at debug level.
// It only fires when the request context carries a logger (set by the request
// logging middleware), so queries outside an HTTP request are not logged.
type QueryTracer struct{}

func (t QueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	return context.WithValue(ctx, queryContext{}, queryContext{
		start: time.Now(),
		sql:   data.SQL,
	})
}

func (t QueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	info, _ := ctx.Value(queryContext{}).(queryContext)
	attrs := []any{
		"sql", info.sql,
		"duration_ms", time.Since(info.start).Milliseconds(),
	}
	if data.Err != nil {
		attrs = append(attrs, "error", data.Err)
		logging.FromContext(ctx).Error("query failed", attrs...)
		return
	}
	attrs = append(attrs, "rows_affected", data.CommandTag.RowsAffected())
	logging.FromContext(ctx).Debug("query", attrs...)
}
