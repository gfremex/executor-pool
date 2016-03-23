package executor_pool

import (
	"testing"
	"strconv"
	"sync"
)

type t_Job_1 struct {
	*BasicJob
	t  *testing.T
	ch chan int64
}

func (job *t_Job_1) Run() {
	var i, sum int64

	for i = 1; i < 101; i++ {
		sum += i
	}

	job.ch <- sum
}

type t_Job_2 struct {
	*BasicJob
	t *testing.T
}

func (job *t_Job_2) Run() {
	panic(-1)
}

func t_NewExecutorPool(nSize, nJobBufSize int) (*ExecutorPool, error) {
	if nSize == 0 && nJobBufSize == 0 {
		return DefaultExecutorPool, nil
	} else {
		return NewExecutorPool(nSize, nJobBufSize)
	}
}

type t_func_NewExecutorPool func(nSize, nJobBufSize int) (*ExecutorPool, error)

func t_test_1(t *testing.T, nSize, nJobBufSize int, nJobs int64, fn t_func_NewExecutorPool) {
	pool, err := fn(nSize, nJobBufSize)

	if err != nil {
		t.Errorf("Can not create new executor pool: %v", err)
	}

	var (
		t_ch chan int64 = make(chan int64, 100)

		t_total int64

		t_wg sync.WaitGroup
	)

	t_wg.Add(1)
	go func() {
		defer t_wg.Done()

		var i int64

		for i = 0; i < nJobs; i++ {
			sum := <-t_ch

			t_total += sum
		}
	}()

	var i int64

	for i = 0; i < nJobs; i++ {
		job := &t_Job_1{
			&BasicJob{
				Id: strconv.FormatInt(i, 10),
			},
			t,
			t_ch,
		}

		pool.Schedule(job)
	}

	t_wg.Wait()

	pool.CloseAndWait()

	t.Logf("Scheduled jobs: %d", pool.Scheduled())

	t.Logf("Total: %d, Sum for per job: %d", t_total, t_total / nJobs)
}

func t_test_2(t *testing.T, nSize, nJobBufSize int, nJobs int64, fn t_func_NewExecutorPool) {
	pool, err := fn(nSize, nJobBufSize)

	if err != nil {
		t.Errorf("Can not create new executor pool: %v", err)
	}

	var i int64

	for i = 0; i < nJobs; i++ {
		job := &t_Job_2{
			&BasicJob{
				Id: strconv.FormatInt(i, 10),
			},
			t,
		}

		pool.Schedule(job)
	}

	pool.CloseAndWait()

	t.Logf("Scheduled jobs: %d", pool.Scheduled())
}

func TestExecutorPool(t *testing.T) {
	var nJobs int64 = 10000000

	t_test_1(t, 20, 100, nJobs, t_NewExecutorPool)

	t_test_1(t, 0, 0, nJobs, t_NewExecutorPool)

	// panic: send on closed channel
	//t_test(t, 0, 0, nJobs, t_NewExecutorPool)

	t_test_2(t, 20, 100, nJobs, t_NewExecutorPool)
}
