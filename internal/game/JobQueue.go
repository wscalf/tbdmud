package game

import "sync"

type JobQueue struct {
	pending sync.Map
	jobs    chan Job
}

func NewJobQueue(maxDepth int) *JobQueue {
	return &JobQueue{
		pending: sync.Map{},
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
			j.pending.Store(job, struct{}{})
			job.SetRequeueHandler(func() {
				j.Requeue(job)
			})
		}
	}
}

func (j *JobQueue) Enqueue(job Job) {
	j.jobs <- job
}

func (j *JobQueue) Requeue(job Job) {
	j.pending.Delete(job)
	j.Enqueue(job)
}
