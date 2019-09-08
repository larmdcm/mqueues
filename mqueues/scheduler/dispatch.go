package scheduler

import (
	"mtypes"
)

type Dispatch struct {
	WorkerChan chan mtypes.Job
}

func (self *Dispatch) WorkeChan() chan mtypes.Job {
	return self.WorkerChan
}

func (self *Dispatch) Run() {
	self.WorkerChan = make(chan mtypes.Job)
}

func (self *Dispatch) Submit(job mtypes.Job) {
	go func() {
		self.WorkerChan <- job
	}()
}
