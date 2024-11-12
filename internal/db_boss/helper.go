package db_boss

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrDuplicateKey = "provided entry already exists"
	ErrNoRows       = "provided entry does not exist"
)

func CheckErrPSQL(data error) error {
	var pgErr *pgconn.PgError
	if errors.As(data, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return errors.New(ErrDuplicateKey)
		case "23503":
			return errors.New("foreign key constraint violation")
		case "42601":
			return errors.New("syntax error in SQL query")
		case "42703":
			return errors.New("invalid column specified in query")
		case "42501":
			return errors.New("insufficient privileges to perform this operation")
		case "08001", "08006":
			return errors.New("database connection error")
		case "40001":
			return errors.New("transaction deadlock or serialization failure; consider retrying")
		case "57014":
			return errors.New("query timed out, try again later")
		case "42804":
			return errors.New("data type mismatch error")
		}
	}
	if errors.Is(data, pgx.ErrNoRows) {
		return errors.New(ErrNoRows)
	}
	return errors.New("an unexpected database error occurred")
}

func CheckErrMySQL(data error) error {
	var mySQLErr *mysql.MySQLError
	if errors.As(data, &mySQLErr) {
		switch mySQLErr.Number {
		case 1062:
			return errors.New(ErrDuplicateKey)
		case 1451, 1452:
			return errors.New("foreign key constraint violation")
		case 1044, 1045:
			return errors.New("insufficient privileges to perform this operation")
		case 2002, 2006:
			return errors.New("database connection error")
		case 1213:
			return errors.New("transaction deadlock or serialization failure; consider retrying")
		case 1406:
			return errors.New("data too long for column")
		}
	}
	if errors.Is(data, sql.ErrNoRows) {
		return errors.New(ErrNoRows)
	}
	return errors.New("an unexpected database error occurred")
}
