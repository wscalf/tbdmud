package game

type JobQueue struct {
	pending map[Job]interface{}
	jobs    chan Job
}

func NewJobQueue(maxDepth int) *JobQueue {
	return &JobQueue{
		pending: make(map[Job]interface{}),
		jobs:    make(chan Job, maxDepth),
	}
}

func (j *JobQueue) Run() {
	for job := range j.jobs {
		job.SetRequeueHandler(func() {
			j.Requeue(job)
		})

		job.Run()
		if !job.IsDone() {
			j.pending[job] = struct{}{}
		}
	}
}

func (j *JobQueue) Enqueue(job Job) {
	j.jobs <- job
}

func (j *JobQueue) Requeue(job Job) {
	j.Enqueue(job)
	delete(j.pending, job)
}
