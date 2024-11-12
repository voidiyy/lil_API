package db_psql

import (
	types "gigaAPI/internal/type"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkerDB(t *testing.T) {
	var err error
	db, err := NewPSQL("../../.env")
	if err != nil {
		t.Error(err)
	}

	t.Run("CreateWorker", func(t *testing.T) {
		worker := CreateTestWorker()
		err = db.WorkerHandler.CreateWorker(worker)
		assert.NoError(t, err, "Failed to create worker")
	})

	t.Run("IfWorkerExists", func(t *testing.T) {
		worker := CreateTestWorker()
		err = db.WorkerHandler.CreateWorker(worker)
		assert.NoError(t, err, "Failed to create worker for existence check")

		err = db.WorkerHandler.IfWorkerExists(worker.ID)
		assert.NoError(t, err, "Error checking if worker exists")
	})

	t.Run("UpdateWorker", func(t *testing.T) {
		worker := CreateTestWorker()
		err = db.WorkerHandler.CreateWorker(worker)
		assert.NoError(t, err, "Failed to create worker for update")
		worker.Password = HashPass(worker.Password)
		worker.Name = genPassORName(14)
		err = db.WorkerHandler.UpdateWorker(worker)
		assert.NoError(t, err, "Failed to update worker")
	})

	t.Run("UpdateWorkerPassword", func(t *testing.T) {
		worker := CreateTestWorker()
		err = db.WorkerHandler.CreateWorker(worker)
		assert.NoError(t, err, "Failed to create worker for password update")

		newpass := "ncdsncdjckskjnkj"
		err = db.WorkerHandler.UpdateWorkerPassword(worker.ID, HashPass(newpass))
		assert.NoError(t, err, "Failed to update worker password")
	})

	t.Run("GetWorker", func(t *testing.T) {
		worker := CreateTestWorker()
		err = db.WorkerHandler.CreateWorker(worker)
		assert.NoError(t, err, "Failed to create worker for retrieval")

		var result = &types.Worker{}
		result, err = db.WorkerHandler.GetByID(worker.ID)
		assert.NoError(t, err, "Failed to retrieve worker by ID")
		assert.Equal(t, worker.ID, result.ID, "Retrieved worker ID should match created worker ID")
	})

	t.Run("SearchWorkers", func(t *testing.T) {
		worker := CreateTestWorker()
		err = db.WorkerHandler.CreateWorker(worker)
		assert.NoError(t, err, "Failed to create worker for search")

		var results []*types.Worker
		results, err = db.WorkerHandler.SearchWorkers(worker.Name)
		assert.NoError(t, err, "Failed to search workers by name")
		assert.True(t, len(results) > 0, "Search should return at least one worker")
	})

	t.Run("GetWorkersByRole", func(t *testing.T) {
		worker := CreateTestWorker()
		err = db.WorkerHandler.CreateWorker(worker)
		assert.NoError(t, err, "Failed to create worker for role filter")

		var results []*types.Worker
		results, err = db.WorkerHandler.GetWorkersByRole("Junior")
		assert.NoError(t, err, "Failed to retrieve workers by role")
		assert.True(t, len(results) > 0, "Should return at least one worker for the given role")
	})

	t.Run("DeactivateWorker", func(t *testing.T) {
		worker := CreateTestWorker()
		err = db.WorkerHandler.CreateWorker(worker)
		assert.NoError(t, err, "Failed to create worker for deactivation")

		err = db.WorkerHandler.DeactivateWorker(worker.ID)
		assert.NoError(t, err, "Failed to deactivate worker")
	})

	t.Run("DeleteWorker", func(t *testing.T) {
		worker := CreateTestWorker()
		err = db.WorkerHandler.CreateWorker(worker)
		assert.NoError(t, err, "Failed to create worker for deletion")

		err = db.WorkerHandler.DeleteWorker(worker.ID)
		assert.NoError(t, err, "Failed to delete worker")
	})

	t.Run("GetCountWorkers", func(t *testing.T) {
		var count int
		count, err = db.WorkerHandler.GetCountWorkers()
		assert.NoError(t, err, "Failed to get worker count")
		assert.True(t, count >= 0, "Worker count should be non-negative")
	})

	t.Run("GetListWorkers", func(t *testing.T) {
		_, err = db.WorkerHandler.GetListWorkers(10, 0)
		assert.NoError(t, err, "Failed to get list of workers")
	})
}
