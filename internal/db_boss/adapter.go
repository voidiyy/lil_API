package db_boss

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBAdapter interface {
	Query(ctx context.Context, query string, args ...interface{}) (UniRows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) UniRow
}

type UniRows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

type UniRow interface {
	Scan(dest ...interface{}) error
}

type UniResult struct {
	commandTag pgconn.CommandTag
}

func (r UniResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("LastInsertId not supported by pgx")
}

func (r UniResult) RowsAffected() (int64, error) {
	return r.commandTag.RowsAffected(), nil
}

func NewUniResult(tag pgconn.CommandTag) sql.Result {
	return UniResult{commandTag: tag}
}
