## executor-pool


Run concurrently with pooled executors (goroutines).

It's caller's responsibility to create a **job** before send it to **executor-pool** to run.

## How to use

### Import package

```Go
import epool "github.com/gfremex/executor-pool"
```

### Create a job

A job is an object which implements **Job** interface.

**Job** interface is defined as:

```Go
type Job interface {
	// Get the Id of the job.
	GetId() string

	// Run the job.
	Run()
}
```

Create a job like this:

```Go
type JobExample struct {
	id string
}

func (job *JobExample) GetId() string {
	return job.id
}

func (job *JobExample) Run() {
	// Do whatever you need
	......
}

var job = &JobExample{"job1"}
```

If panics within **job.Run()**, an error message will be logged by log.Printf(...).

### Create an executor pool

You can pass two parameters when creating:

- size (*Number of executors*)
- jobBufSize (*Job buffer size*)

For example:

```Go
// Create a new executor pool with:
// 10 executors
// 100 buffered jobs
pool, err := epool.NewExecutorPool(10, 100)
```

### Run a job

Call **pool.Schedule(job)**

The caller goroutine will be blocked until the job is accepted.

### Get the number of run jobs

Call **pool.Scheduled()**

### Close an executor pool

If no more jobs need to run, you can close the pool.

#### pool.Close()

Just notify closing.

#### pool.CloseAndWait()

Notify closing and wait until the pool is closed.

### Shorcuts for instant use

A builtin executor pool is provided for instant use.

- epool.Schedule(job)
- epool.Scheduled()
- epool.Close()
- epool.CloseAndWait()

***Note:***

Do not call ***epool.Close()*** or ***epool.CloseAndWait()*** if anyone else is also using the builtin executor pool.

Generally ***epool.Close()*** and ***epool.CloseAndWait()*** should be called in a deferred function of main goroutine.

#### Complete examples

[https://github.com/gfremex/executor-pool/tree/master/examples](https://github.com/gfremex/executor-pool/tree/master/examples)
