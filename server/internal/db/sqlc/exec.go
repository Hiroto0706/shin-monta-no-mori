package db

import (
	"context"
	"database/sql"
)

func (q *Queries) ExecQuery(ctx context.Context, query string) (sql.Result, error) {
	return q.db.ExecContext(ctx, query)
}
