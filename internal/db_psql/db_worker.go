package db_psql

import (
	"context"
	"errors"
	"fmt"
	"gigaAPI/internal/type"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var (
	_ workerDB = &workerHandler{}
)

type workerDB interface {
	CreateWorker(wrk *types.Worker) error

	DeactivateWorker(id uuid.UUID) error
	DeleteWorker(id uuid.UUID) error

	LoginWorker(email, pass string) (*types.Worker, error)

	SearchWorkers(query string) ([]*types.Worker, error)

	GetWorkersByRole(role string) ([]*types.Worker, error)
	GetWorkerByEmail(email string) (*types.Worker, error)
	GetListWorkers(limit, offset int) ([]*types.Worker, error)
	GetCountWorkers() (int, error)
	GetByID(id uuid.UUID) (*types.Worker, error)

	IfWorkerExists(id uuid.UUID) error

	UpdateWorker(wrk *types.Worker) error
	UpdateWorkerPassword(id uuid.UUID, newPassword string) error
}
type workerHandler struct {
	Context context.Context
	*pgxpool.Pool
}

func (p *workerHandler) GetWorkersByRole(role string) ([]*types.Worker, error) {
	query := `SELECT id, name, email, role FROM worker WHERE role = $1`
	rows, err := p.Query(p.Context, query, role)
	if err != nil {
		return nil, fmt.Errorf("error fetching workers by role '%s': %w", role, err)
	}
	defer rows.Close()

	var workers []*types.Worker
	for rows.Next() {
		var worker types.Worker
		if err := rows.Scan(&worker.ID, &worker.Name, &worker.Email, &worker.Role); err != nil {
			return nil, fmt.Errorf("error scanning worker by role '%s': %w", role, err)
		}
		workers = append(workers, &worker)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows for role '%s': %w", role, err)
	}
	return workers, nil
}

func (p *workerHandler) GetWorkerByEmail(email string) (*types.Worker, error) {
	var worker types.Worker
	query := `SELECT id, name, email, role FROM worker WHERE email=$1`
	row := p.QueryRow(p.Context, query, email)
	if err := row.Scan(&worker.ID, &worker.Name, &worker.Email, &worker.Role); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no worker found with email '%s'", email)
		}
		return nil, fmt.Errorf("error fetching worker by email '%s': %w", email, err)
	}
	return &worker, nil
}

func (p *workerHandler) GetListWorkers(limit, offset int) ([]*types.Worker, error) {
	query := `SELECT id, name, email, role FROM worker LIMIT $1 OFFSET $2`
	rows, err := p.Query(p.Context, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing workers with limit %d and offset %d: %w", limit, offset, err)
	}
	defer rows.Close()

	var workers []*types.Worker
	for rows.Next() {
		var worker types.Worker
		if err := rows.Scan(&worker.ID, &worker.Name, &worker.Email, &worker.Role); err != nil {
			return nil, fmt.Errorf("error scanning worker in list: %w", err)
		}
		workers = append(workers, &worker)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through worker list rows: %w", err)
	}
	return workers, nil
}

func (p *workerHandler) GetCountWorkers() (int, error) {
	var count int
	query := `SELECT count(*) FROM worker`
	row := p.QueryRow(p.Context, query)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("error getting worker count: %w", err)
	}
	return count, nil
}

func (p *workerHandler) IfWorkerExists(id uuid.UUID) error {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM worker WHERE id=$1)`
	row := p.QueryRow(p.Context, query, id)
	if err := row.Scan(&exists); err != nil {
		return fmt.Errorf("error checking existence of worker with id %s: %w", id, err)
	}
	if exists {
		return nil
	}
	return fmt.Errorf("worker not exist by id: %s ", id)
}

func (p *workerHandler) UpdateWorkerPassword(id uuid.UUID, newPassword string) error {
	query := `UPDATE worker SET password = $1 WHERE id = $2`
	cmd, err := p.Exec(p.Context, query, newPassword, id)
	if err != nil {
		return fmt.Errorf("error updating password for worker with id %s: %w", id, err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected when updating password for worker with id %s", id)
	}
	return nil
}

var ErrDuplicateEntry = errors.New("entry already exists")

func (p *workerHandler) GetByID(id uuid.UUID) (*types.Worker, error) {
	w := &types.Worker{}
	query := `SELECT id, name, password,email, role, created_at FROM worker WHERE id = $1`
	row := p.QueryRow(p.Context, query, id)
	if err := row.Scan(&w.ID, &w.Name, &w.Password, &w.Email, &w.Role, &w.CreatedAt); err != nil {
		log.Printf("Error retrieving worker by ID %s: %v", id, err)
		return nil, err
	}
	return w, nil
}

func (p *workerHandler) CreateWorker(wrk *types.Worker) error {
	query := `INSERT INTO worker (id, name, password, email, role) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.Exec(p.Context, query, wrk.ID, wrk.Name, wrk.Password, wrk.Email, wrk.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Printf("Duplicate worker entry: %v", err)
			return fmt.Errorf("%w: a worker with this email or name already exists", ErrDuplicateEntry)
		}
		log.Printf("Error creating worker: %v", err)
		return fmt.Errorf("error creating worker: %v", err)
	}

	log.Printf("Worker created successfully: %s", wrk.Name)
	return nil
}

func (p *workerHandler) DeactivateWorker(id uuid.UUID) error {
	query := `UPDATE worker SET isactive = false WHERE id = $1`
	result, err := p.Exec(p.Context, query, id)
	if err != nil {
		log.Printf("Error deactivating worker %s: %v", id, err)
		return fmt.Errorf("error deactivating worker: %v", err)
	}
	if result.RowsAffected() == 0 {
		log.Printf("No worker found with ID %s for deactivation", id)
		return fmt.Errorf("no worker found with the specified ID: %s", id)
	}
	log.Printf("Worker deactivated successfully: %s", id)
	return nil
}

func (p *workerHandler) UpdateWorker(wrk *types.Worker) error {
	query := `UPDATE worker SET name = $2, password = $3, email = $4, role = $5 WHERE id = $1`
	result, err := p.Exec(p.Context, query, wrk.ID, wrk.Name, wrk.Password, wrk.Email, wrk.Role)
	if err != nil {
		log.Printf("Error updating worker %s: %v", wrk.ID, err)
		return fmt.Errorf("error updating worker: %v", err)
	}
	if result.RowsAffected() == 0 {
		log.Printf("No rows affected while updating worker %s", wrk.ID)
		return fmt.Errorf("no worker found with the specified ID: %s", wrk.ID)
	}
	log.Printf("Worker updated successfully: %s", wrk.Name)
	return nil
}

func (p *workerHandler) LoginWorker(email, pass string) (*types.Worker, error) {
	w := &types.Worker{}
	query := `SELECT id, name, password, email, role, created_at FROM worker WHERE email = $1`
	row := p.QueryRow(p.Context, query, email)
	err := row.Scan(&w.ID, &w.Name, &w.Password, &w.Email, &w.Role, &w.CreatedAt)
	if err != nil {
		log.Printf("Error scanning worker for login by Email %s: %v", email, err)
		return nil, fmt.Errorf("error scanning worker: %v", err)
	}
	fmt.Println("input pass: ", pass)
	fmt.Println("db pass: ", w.Password)
	if !HashCompare(w.Password, pass) {
		log.Println("Login pass mismatch")
		return nil, fmt.Errorf("error login worker: %v", fmt.Errorf("password mismatch"))
	}
	return w, nil
}

func (p *workerHandler) DeleteWorker(id uuid.UUID) error {
	result, err := p.Exec(p.Context, "DELETE FROM worker WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting worker %s: %v", id, err)
		return fmt.Errorf("error deleting worker: %v", err)
	}
	if result.RowsAffected() == 0 {
		log.Printf("No worker found with ID %s for deletion", id)
		return fmt.Errorf("no worker found with the specified ID: %s", id)
	}
	log.Printf("Worker deleted successfully: %s", id)
	return nil
}

func (p *workerHandler) SearchWorkers(query string) ([]*types.Worker, error) {
	searchQuery := `SELECT id, name, email, role FROM worker WHERE name LIKE '%' || $1 || '%'`
	rows, err := p.Query(p.Context, searchQuery, query)
	if err != nil {
		log.Printf("Error searching workers with query '%s': %v", query, err)
		return nil, fmt.Errorf("error searching workers: %v", err)
	}
	defer rows.Close()

	var workers []*types.Worker
	for rows.Next() {
		var wrk types.Worker
		if err = rows.Scan(&wrk.ID, &wrk.Name, &wrk.Email, &wrk.Role); err != nil {
			log.Printf("Error scanning worker result: %v", err)
			return nil, fmt.Errorf("error scanning worker: %v", err)
		}
		workers = append(workers, &wrk)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error after search: %v", err)
		return nil, fmt.Errorf("error iterating worker rows: %v", err)
	}

	log.Printf("Search completed successfully with %d results", len(workers))
	return workers, nil
}
