package mqueues

import (
	"mtypes"
	"encoding/json"
)

type Job struct {
	JobData *mtypes.Job
	queue *Queue
}

func (self *Job) Delete () error {
	jobK,err := self.GetJobDataJson()
	if err != nil {
		return err
	}
	return self.queue.Conn.Delete(self.queue.Conn.GetQueueExcutingName(),jobK)
}

func (self *Job) Release (delay int64) error {
	err := self.Delete()
	if err != nil {
		return err
	}
	self.JobData.AttemptsCount++
	return self.queue.Conn.Later(self.queue.Conn.GetQueueName(),delay,self.JobData)
}

func (self *Job) GetJobDataJson () ([]byte,error) {
	return json.Marshal(self.JobData)
}

func (self *Job) Attempts () int64 {
	return self.JobData.AttemptsCount
}