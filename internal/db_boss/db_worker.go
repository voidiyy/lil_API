package db_boss

import (
	"context"
	"fmt"
	"gigaAPI/internal/logger"
	types "gigaAPI/internal/type"
	"github.com/google/uuid"
	"time"
)

// for single result (SELECT) - QueryRow
// for multiple results (SELECT) - Query
// for changing data (INSERT, DELETE, UPDATE) -Exec

type workerHandler interface {
	CreateW(ctx context.Context, wrk *types.Worker) error
	DeleteW(ctx context.Context, id uuid.UUID) error
	UpdateW(ctx context.Context, wrk *types.Worker) (*types.Worker, error)
	GetW(ctx context.Context, id uuid.UUID) (*types.Worker, error)
	GetWByEmail(ctx context.Context, email string) (*types.Worker, error)
}

type WorkerH struct {
	db *database
	lg *logger.Logger
}

func NewWorkerH(db *database) *WorkerH {
	return &WorkerH{
		db: db,
		lg: logger.NewLogger(),
	}
}

func (w *WorkerH) GetWByEmail(ctx context.Context, email string) (*types.Worker, error) {
	var op = tracer("GetWByEmail")
	worker := &types.Worker{}

	switch t := w.db.adapter.(type) {
	case *PostgresAdapter:
		q := `select id, name, password, email, role, isactive, created_at from worker where email=$1`
		row := t.QueryRow(ctx, q, email)
		err := row.Scan(&worker.ID, &worker.Name, &worker.Password, &worker.Email, &worker.Role, &worker.IsActive, &worker.CreatedAt)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrPSQL(err)
			return nil, e
		}
		w.lg.LogDBQuery(q)
		return worker, nil
	case *MySQLAdapter:
		q := `select id, name, password, email, role, isactive, created_at from worker where email=?`
		row := t.QueryRow(ctx, q, email)
		var createdAt []byte
		err := row.Scan(&worker.ID, &worker.Name, &worker.Password, &worker.Email, &worker.Role, &worker.IsActive, &createdAt)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrMySQL(err)
			return nil, e
		}
		worker.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			w.lg.LogDBError(err, op)
			return nil, err
		}

		w.lg.LogDBQuery(q)
		return worker, nil

	default:
		err := fmt.Errorf("error db type accertion")
		w.lg.LogDBError(err, op)
		return nil, err
	}
}

func (w *WorkerH) GetW(ctx context.Context, id uuid.UUID) (*types.Worker, error) {
	var op = tracer("GetW")
	worker := &types.Worker{}

	switch t := w.db.adapter.(type) {
	case *PostgresAdapter:
		q := `select id, name, password, email, role, isactive, created_at from worker where id=$1`
		row := t.QueryRow(ctx, q, id)
		err := row.Scan(&worker.ID, &worker.Name, &worker.Password, &worker.Email, &worker.Role, &worker.IsActive, &worker.CreatedAt)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrPSQL(err)
			return nil, e
		}
		w.lg.LogDBQuery(q)
		return worker, nil
	case *MySQLAdapter:
		q := `select id, name, password, email, role, isactive, created_at from worker where id=?`
		row := t.QueryRow(ctx, q, id)
		var createdAt []byte
		err := row.Scan(&worker.ID, &worker.Name, &worker.Password, &worker.Email, &worker.Role, &worker.IsActive, &createdAt)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrMySQL(err)
			return nil, e
		}
		worker.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			w.lg.LogDBError(err, op)
			return nil, err
		}

		w.lg.LogDBQuery(q)
		return worker, nil

	default:
		err := fmt.Errorf("error db type accertion")
		w.lg.LogDBError(err, op)
		return nil, err
	}
}

func (w *WorkerH) UpdateW(ctx context.Context, wrk *types.Worker) (*types.Worker, error) {
	var op = tracer("UpdateW")
	switch t := w.db.adapter.(type) {
	case *PostgresAdapter:
		q := `update worker set name=$2, password=$3, email=$4, role=$5, isactive=$6 where id=$1`
		_, err := t.Exec(ctx, q, wrk.ID, wrk.Name, wrk.Password, wrk.Email, wrk.Role, true)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrPSQL(err)
			return nil, e
		}
		w.lg.LogDBQuery(q)
		return w.GetW(ctx, wrk.ID)
	case *MySQLAdapter:
		q := `update worker set name=?, password=?, email=?, role=?, isactive=true where id=?`
		_, err := t.Exec(ctx, q, wrk.Name, wrk.Password, wrk.Email, wrk.Role, wrk.ID)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrMySQL(err)
			return nil, e
		}
		w.lg.LogDBQuery(q)
		return w.GetW(ctx, wrk.ID)
	default:
		err := fmt.Errorf("error db type accertion")
		w.lg.LogDBError(err, op)
		return nil, err
	}
}

func (w *WorkerH) CreateW(ctx context.Context, wrk *types.Worker) error {
	var op = tracer("CreateW")

	switch t := w.db.adapter.(type) {
	case *PostgresAdapter:
		q := `insert into worker (id, name, password, email, role, isactive) values ($1, $2,$3,$4,$5,$6)`
		_, err := t.Exec(ctx, q, &wrk.ID, &wrk.Name, &wrk.Password, &wrk.Email, &wrk.Role, true)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrPSQL(err)
			return e
		}
		w.lg.LogDBQuery(q)
		return nil
	case *MySQLAdapter:
		q := `insert into worker (id, name, password, email, role, isactive) values (?, ?,?,?,?,?)`
		_, err := t.Exec(ctx, q, &wrk.ID, &wrk.Name, &wrk.Password, &wrk.Email, &wrk.Role, true)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrPSQL(err)
			return e
		}
		w.lg.LogDBQuery(q)
		return nil
	default:
		err := fmt.Errorf("error db type accertion")
		w.lg.LogDBError(err, op)
		return err
	}
}

func (w *WorkerH) DeleteW(ctx context.Context, id uuid.UUID) error {
	var op = tracer("DeleteW")

	switch t := w.db.adapter.(type) {
	case *PostgresAdapter:
		q := `delete from worker where id = $1`
		_, err := t.Exec(ctx, q, id)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrPSQL(err)
			return e
		}
		w.lg.LogDBQuery(q)
		return nil
	case *MySQLAdapter:
		q := `delete from worker where id = ?`
		_, err := t.Exec(ctx, q, id)
		if err != nil {
			w.lg.LogDBError(err, op)
			e := CheckErrMySQL(err)
			return e
		}
		w.lg.LogDBQuery(q)
		return nil
	default:
		err := fmt.Errorf("error db type accertion")
		w.lg.LogDBError(err, op)
		return err
	}
}
