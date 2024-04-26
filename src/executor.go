package executer


import (
	"fmt"
	"sync"
)

// ParsedStmtInterface represents the parsed SQL statement.
type ParsedStmtInterface interface {
	GetQueryType() QueryType
	GetTables() []string
	GetColumns() []string
	GetConditions() []interface{}
	GetValues() []string
}

// QueryType represents the type of SQL query.
type QueryType int

const (
	SELECT QueryType = iota
	INSERT
	UPDATE
	DELETE
)

// Worker represents a worker that processes SQL queries.
type Worker struct {
	id       int
	busy     bool
	caches   map[string]*Cache
	storage  *StorageEngine
}
// NewWorker creates a new Worker instance.
func NewWorker(id int, caches map[string]*Cache, storage *StorageEngine) *Worker {
	return &Worker{
		id: id,
		caches: caches, 
		storage: storage,
	}
}

// QueryExecutor represents the query executor.
type QueryExecutor struct {
	workers   []*Worker
	workerMux sync.RWMutex
	caches    map[string]*Cache
	storage   *StorageEngine
}

// NewQueryExecutor creates a new QueryExecutor instance.
func NewQueryExecutor(numWorkers int, caches map[string]*Cache, storage *StorageEngine) *QueryExecutor {
	executor := &QueryExecutor{
		caches:  caches,
		storage: storage,
	}
	for i := 0; i < numWorkers; i++ {
		executor.createWorker()
	}
	return executor
}

// ExecuteStatement assigns the parsed statement to a worker for execution.
func (qe *QueryExecutor) ExecuteStatement(stmt ParsedStmtInterface) {
	// Get an available worker or create a new one if all are busy
	worker := qe.getAvailableWorker()
	worker.ExecuteQuery(stmt)
}

// ExecuteQuery executes the SQL query assigned to the worker.
func (w *Worker) ExecuteQuery(stmt ParsedStmtInterface) {
	w.busy = true
	defer func() { w.busy = false }()

	switch stmt.GetQueryType() {
	case SELECT:
		tables := stmt.GetTables()
		for _, table := range tables {
			cache, ok := w.caches[table]
			if !ok {
				fmt.Printf("Cache %s not found\n", table)
				continue
			}
			if !cache.Contains(table) {
				fmt.Printf("Cache miss for table %s\n", table)
				// Retrieve data from storage
				         // Cache the data
			}
		}
		// Implement SELECT query execution logic
	case INSERT, UPDATE, DELETE:
		// Implement write query execution logic
		// Write to cache first and then sync with storage
		// For simplicity, assuming all write operations are done on single tables
		table := stmt.GetTables()[0]
		cache, ok := w.caches[table]
		if !ok {
			fmt.Printf("Cache %s not found\n", table)
			return
		}
		cache.Invalidate(table) // Invalidate cache for the table
		// Write to cache
		// Write to storage
	}
}

// getAvailableWorker returns an available worker or creates a new one if all are busy.
func (qe *QueryExecutor) getAvailableWorker() *Worker {
	qe.workerMux.Lock()
	defer qe.workerMux.Unlock()

	for _, worker := range qe.workers {
		// If worker is not busy, return it
		if !worker.busy {
			return worker
		}
	}

	// If all workers are busy, create a new one
	return qe.createWorker()
}

// createWorker creates a new worker with access to the same caches and storage engine as existing workers.
func (qe *QueryExecutor) createWorker() *Worker {
	newWorker := NewWorker(len(qe.workers), qe.caches, qe.storage)
	qe.workers = append(qe.workers, newWorker)
	return newWorker
}

