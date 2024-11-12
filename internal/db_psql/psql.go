package db_psql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"os"
)

type HandlerDB struct {
	WorkerHandler workerDB
}

func NewPSQL(path string) (*HandlerDB, error) {
	err := godotenv.Load(path)
	if err != nil {
		fmt.Println("env trouble: ", err)
		return nil, err
	}

	dburl := os.Getenv("DB_URL")
	dbconfig, er := pgxpool.ParseConfig(dburl)
	if er != nil {
		fmt.Println("db url trouble: ", err)
		return nil, er
	}
	ctx := context.Background()

	pool, e := pgxpool.NewWithConfig(ctx, dbconfig)
	if e != nil {
		fmt.Println("db config trouble: ", e)
		return nil, e
	}

	if pool == nil {
		return nil, fmt.Errorf("pool is nil ! ")
	}
	err = pool.Ping(ctx)
	if err != nil {
		fmt.Println("db ping error: ", err)
		return nil, err
	}

	db := &HandlerDB{
		WorkerHandler: &workerHandler{
			Context: ctx,
			Pool:    pool,
		},
	}

	return db, nil
}
