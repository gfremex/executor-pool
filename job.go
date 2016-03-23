package executor_pool

type Job interface {
	// Get the Id of the job.
	GetId() string

	// Run the job.
	Run()
}

type BasicJob struct {
	// Id of the job
	Id string
}

func (job *BasicJob) GetId() string {
	return job.Id
}

type RunFunc func()
