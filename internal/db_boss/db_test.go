package db_boss

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gigaAPI/config"
	"gigaAPI/internal/db_psql"
	"testing"
)

func TestInitDB_Postgres(t *testing.T) {
	p := "../../postgres.yaml"

	conf := config.InitConfig(p)

	db, err := InitDB(conf.DbConf.DbType, conf.DbConf.DbURL)
	if err != nil {
		t.Error(err)
	}

	t.Run("psql_CreateW", func(t *testing.T) {
		worker := db_psql.CreateTestWorker()

		err = db.Worker.CreateW(context.Background(), worker)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("psql_GetW", func(t *testing.T) {
		worker := db_psql.CreateTestWorker()

		err = db.Worker.CreateW(context.Background(), worker)
		if err != nil {
			t.Error(err)
		}

		w, er := db.Worker.GetW(context.Background(), worker.ID)
		if er != nil {
			t.Error(er)
		}

		if worker.Name != w.Name {
			t.Errorf("worker name mismatch: was: %s get: %s", worker.Name, w.Name)
		}
	})
	t.Run("psql_DeleteW", func(t *testing.T) {
		worker := db_psql.CreateTestWorker()

		err = db.Worker.CreateW(context.Background(), worker)
		if err != nil {
			t.Error(err)
		}

		err = db.Worker.DeleteW(context.Background(), worker.ID)
		if err != nil {
			t.Error(err)
		}

		_, er := db.Worker.GetW(context.Background(), worker.ID)
		if !errors.Is(er, sql.ErrNoRows) {
			fmt.Println("ERROR", er)
		}
	})
	t.Run("psql_UpdateW", func(t *testing.T) {
		worker := db_psql.CreateTestWorker()

		err = db.Worker.CreateW(context.Background(), worker)
		if err != nil {
			t.Error(err)
		}

		newWorker := db_psql.CreateTestWorker()
		newWorker.ID = worker.ID

		w, er := db.Worker.UpdateWorker(context.Background(), newWorker)
		if er != nil {
			t.Error(er)
		}

		if w.Name != newWorker.Name {
			t.Errorf("name mismatch: newName: %s returned name: %s", newWorker.Name, w.Name)
		}
	})
}

func TestInitDB_MySQL(t *testing.T) {

	p := "../../mysql.yaml"

	conf := config.InitConfig(p)

	db, err := InitDB(conf.DbConf.DbType, conf.DbConf.DbURL)
	if err != nil {
		t.Error(err)
	}
	t.Run("psql_CreateW", func(t *testing.T) {
		worker := db_psql.CreateTestWorker()

		err = db.Worker.CreateW(context.Background(), worker)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("psql_GetW", func(t *testing.T) {
		worker := db_psql.CreateTestWorker()

		err = db.Worker.CreateW(context.Background(), worker)
		if err != nil {
			t.Error(err)
		}

		w, er := db.Worker.GetW(context.Background(), worker.ID)
		if er != nil {
			t.Error(er)
		}

		if worker.Name != w.Name {
			t.Errorf("worker name mismatch: was: %s get: %s", worker.Name, w.Name)
		}
	})
	t.Run("psql_DeleteW", func(t *testing.T) {
		worker := db_psql.CreateTestWorker()

		err = db.Worker.CreateW(context.Background(), worker)
		if err != nil {
			t.Error(err)
		}

		err = db.Worker.DeleteW(context.Background(), worker.ID)
		if err != nil {
			t.Error(err)
		}

		_, er := db.Worker.GetW(context.Background(), worker.ID)
		if !errors.Is(er, sql.ErrNoRows) {
			fmt.Println("ERROR", er)
		}
	})
	t.Run("psql_UpdateW", func(t *testing.T) {
		worker := db_psql.CreateTestWorker()

		err = db.Worker.CreateW(context.Background(), worker)
		if err != nil {
			t.Error(err)
		}

		newWorker := db_psql.CreateTestWorker()
		newWorker.ID = worker.ID

		w, er := db.Worker.UpdateWorker(context.Background(), newWorker)
		if er != nil {
			t.Error(er)
		}

		if w.Name != newWorker.Name {
			t.Errorf("name mismatch: newName: %s returned name: %s", newWorker.Name, w.Name)
		}
	})

}
