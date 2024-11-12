package db_boss

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAdapter struct {
	pool *pgxpool.Pool
}

func (p *PostgresAdapter) Query(ctx context.Context, query string, args ...interface{}) (UniRows, error) {
	return p.pool.Query(ctx, query, args...)
}

func (p *PostgresAdapter) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ct, err := p.pool.Exec(ctx, query, args...)
	return NewUniResult(ct), err
}

func (p *PostgresAdapter) QueryRow(ctx context.Context, query string, args ...interface{}) UniRow {
	return p.pool.QueryRow(ctx, query, args...)
}
