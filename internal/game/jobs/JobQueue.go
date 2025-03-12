package jobs

type JobQueue struct {
	jobs chan Job
}

func NewJobQueue(maxDepth int) *JobQueue {
	return &JobQueue{
		jobs: make(chan Job, maxDepth),
	}
}

func (j *JobQueue) Run() {
	for job := range j.jobs {
		job.Run()
	}
}

func (j *JobQueue) Enqueue(job Job) {
	j.jobs <- job
}
