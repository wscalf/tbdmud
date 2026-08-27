package game

type BaseJob struct {
	done    bool
	handler func()
}

func NewBaseJob() BaseJob {
	return BaseJob{
		done: false,
	}
}

func (j *BaseJob) IsDone() bool {
	return j.done
}

func (j *BaseJob) SetRequeueHandler(handler func()) {
	j.handler = handler
}

type Job interface {
	IsDone() bool
	SetRequeueHandler(handler func())
	Run()
}
