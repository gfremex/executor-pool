package main

import (
	epool "github.com/gfremex/executor-pool"
	"strconv"
	"sync"
	"log"
)

var (
	nJobs = 1000
	ch = make(chan int)
	total int
	wg sync.WaitGroup
)

type JobExample struct {
	id string
	ch chan int
}

func (job *JobExample) GetId() string {
	return job.id
}

func (job *JobExample) Run() {
	sum := 0

	for i := 1; i < 101; i++ {
		sum += i
	}

	job.ch <- sum
}

func useCustomizedPool() {
	log.Println("Use customized pool...")

	total = 0

	// Create a new executor pool with:
	// 10 executors
	// 100 buffered jobs
	pool, err := epool.NewExecutorPool(10, 100)

	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < nJobs; i++ {
			sum := <-ch

			total += sum
		}
	}()

	for i := 0; i < nJobs; i++ {
		job := &JobExample{
			strconv.Itoa(i),
			ch,
		}

		pool.Schedule(job)
	}

	// Wait for the result
	wg.Wait()

	// Close the pool
	pool.CloseAndWait()

	log.Printf("Total: %d, Scheduled jobs: %d, Sum for per job: %d\n", total, pool.Scheduled(), total / pool.Scheduled())
}

func useDefaultPool() {
	log.Println("Use default pool...")

	total = 0

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < nJobs; i++ {
			sum := <-ch

			total += sum
		}
	}()

	for i := 0; i < nJobs; i++ {
		job := &JobExample{
			strconv.Itoa(i),
			ch,
		}

		epool.Schedule(job)
	}

	// Wait for the result
	wg.Wait()

	// Close the pool
	epool.CloseAndWait()

	log.Printf("Total: %d, Scheduled jobs: %d, Sum for per job: %d\n", total, epool.Scheduled(), total / epool.Scheduled())
}

func main() {
	useCustomizedPool()

	log.Println()

	useDefaultPool()
}
