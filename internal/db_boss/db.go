package db_boss

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	adapter DBAdapter
}

func newDatabase(a DBAdapter) *database {
	return &database{
		adapter: a,
	}
}

type HandlerDB struct {
	Worker workerHandler
}

func NewHandlerDB(db *database) *HandlerDB {
	return &HandlerDB{
		Worker: NewWorkerH(db),
	}
}

const path = "internal/db_boss/db."

func tracer(s string) string {
	return path + s
}

func InitDB(dbtype, dburl string) (*HandlerDB, error) {
	var op = tracer(".InitDb")

	switch dbtype {
	case "postgres":
		pool, er := initPgx(dburl)
		if er != nil {
			fmt.Println(op, ":", er)
			return nil, er
		}
		db := newDatabase(&PostgresAdapter{
			pool: pool,
		})
		return NewHandlerDB(db), nil
	case "mysql":
		mysql, er := initMySQL(dburl)
		if er != nil {
			fmt.Println(op, ":", er)
			return nil, er
		}
		db := newDatabase(&MySQLAdapter{
			db: mysql,
		})
		return NewHandlerDB(db), nil
	default:
		err := fmt.Errorf("invalid db type, use postgres or mysql")
		fmt.Println(op, ":", err)
		return nil, err
	}
}
func initMySQL(dburl string) (*sql.DB, error) {
	var (
		op = tracer(".initMySQL")
	)

	db, err := sql.Open("mysql", dburl)
	if err != nil {
		fmt.Println(op, ":", err)
		return nil, err
	}
	if db == nil {
		e := fmt.Errorf("mysql connection is nil")
		fmt.Println(op, ":", e)
		return nil, e
	}

	return db, nil
}

func initPgx(dburl string) (*pgxpool.Pool, error) {
	var (
		op = tracer(".initPgx")
	)

	config, err := pgxpool.ParseConfig(dburl)
	if err != nil {
		fmt.Println(op, ":", err)
		return nil, err
	}

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Println(op, ":", err)
		return nil, err
	}
	if pool == nil {
		e := fmt.Errorf("connection pool is nil")
		fmt.Println(op, ":", e)
		return nil, e
	}

	err = pool.Ping(ctx)
	if err != nil {
		fmt.Println(op, ":", err)
	}

	return pool, nil
}
