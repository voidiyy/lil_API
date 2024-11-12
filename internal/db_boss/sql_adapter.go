package db_boss

import (
	"context"
	"database/sql"
)

type MySQLAdapter struct {
	db *sql.DB
}

func (m *MySQLAdapter) Query(ctx context.Context, query string, args ...interface{}) (UniRows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m *MySQLAdapter) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *MySQLAdapter) QueryRow(ctx context.Context, query string, args ...interface{}) UniRow {
	return m.db.QueryRowContext(ctx, query, args...)
}
