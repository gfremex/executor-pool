package executor_pool

import (
	"log"
	"errors"
	"sync"
)

var (
// Default ExecutorPool
	DefaultExecutorPool, _ = NewExecutorPool(50, 1000)
)

// Call Job.Run() safely.
// Panics will be recovered so that the executor goroutine does not end.
func safeRun(job Job) RunFunc {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Error on running job [%s]: %v\n", job.GetId(), err)
			}
		}()

		job.Run()
	}
}

type executor struct {
	jobChan chan Job
	master  *ExecutorPool
}

// Call Job.Run() safely or return on close.
func (exec *executor) do() {
	defer exec.master.wg.Done()

	for {
		if job, ok := <-exec.jobChan; !ok {
			// exec.jobChan is closed
			// This executor needs to be closed

			return
		} else {
			// Call Job.Run() safely
			safeRun(job)()

			// Make self idle
			exec.master.idleExecutors <- exec
		}
	}
}

// Executor pool run Jobs concurrently with reusable executors.
type ExecutorPool struct {
	// Number of executors
	size          int

	// Job buffer size
	jobBufSize    int

	// Channel of Jobs
	jobs          chan Job

	// Channel of idle executors
	idleExecutors chan *executor

	// Number of scheduled jobs
	scheduled     int

	wg            sync.WaitGroup
}

// Create and start a new ExecutorPool.
// size and jobBufSize can not be negative or zero.
func NewExecutorPool(nSize, nJobBufSize int) (*ExecutorPool, error) {
	if nSize <= 0 {
		return nil, errors.New("Size of ExecutorPool can not be negative or zero")
	}

	if nJobBufSize <= 0 {
		return nil, errors.New("Job buffer size of ExecutorPool can not be negative or zero")
	}

	executorPool := &ExecutorPool{
		size: nSize,
		jobBufSize: nJobBufSize,
		jobs: make(chan Job, nJobBufSize),
		idleExecutors: make(chan *executor, nSize),
	}

	// Create and start executors.
	for i := 0; i < nSize; i++ {
		exec := &executor{
			jobChan: make(chan Job),
			master: executorPool,
		}

		// Make exec idle
		executorPool.idleExecutors <- exec

		// Start exec
		executorPool.wg.Add(1)
		go exec.do()
	}

	// Start executorPool
	executorPool.wg.Add(1)
	go dispatch(executorPool)

	return executorPool, nil
}

// Schedule running of a Job.
// If the job is nil, that means to close the ExecutorPool.
func (executorPool *ExecutorPool) Schedule(job Job) {
	if job == nil {
		executorPool.Close()
	} else {
		executorPool.jobs <- job

		executorPool.scheduled++
	}
}

// Get the number of scheduled jobs.
func (executorPool *ExecutorPool) Scheduled() int {
	return executorPool.scheduled
}

// Close the ExecutorPool and return immediately.
func (executorPool *ExecutorPool) Close() {
	close(executorPool.jobs)
}

// Close the ExecutorPool and wait until it finished.
func (executorPool *ExecutorPool) CloseAndWait() {
	executorPool.Close()

	executorPool.wg.Wait()
}

func dispatch(executorPool *ExecutorPool) {
	defer executorPool.wg.Done()

	for {
		if job, ok := <-executorPool.jobs; !ok {
			// executorPool.jobs is closed
			// This pool needs to be closed

			for i := 0; i < executorPool.size; i++ {
				executor := <-executorPool.idleExecutors
				close(executor.jobChan)
			}

			return
		} else {
			// Forward the job to an idle executor

			executor := <-executorPool.idleExecutors

			executor.jobChan <- job
		}
	}
}

func Schedule(job Job) {
	DefaultExecutorPool.Schedule(job)
}

func Scheduled() int {
	return DefaultExecutorPool.Scheduled()
}

func Close() {
	DefaultExecutorPool.Close()
}

func CloseAndWait() {
	DefaultExecutorPool.CloseAndWait()
}
