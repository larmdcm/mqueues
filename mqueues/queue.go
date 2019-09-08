package mqueues

import (
	"mtypes"
	"strconv"
	"time"
	"fmt"
	"sync"
)

var lock sync.Mutex

type Queue struct {
	Name string
	Conn mtypes.Connector
	Scheduler mtypes.Scheduler
	WorkerCount int
	Handles map[string]mtypes.Handler
}

func New (config mtypes.Config) (*Queue,error){
	err :=  config.Connector.Connect(config)
	if err != nil {
		return &Queue{},err
	}
	config.Connector.SetQueueName(config.Name)
	return &Queue{
		Name: config.Name,
		Conn:  config.Connector,
		Scheduler: config.Scheduler,
		WorkerCount: config.WorkerCount,
		Handles: config.Handlers,
	},nil
}

func (self *Queue) Run () {
	self.Scheduler.Run()

	for i := 0; i < self.WorkerCount; i++ {
		createWorker(self.Scheduler.WorkeChan(),self)
	}

	for {
		lock.Lock()
		job,err := self.Pop()
		lock.Unlock()
		if err != nil{
			state,_ := strconv.Atoi(err.Error())

			if state == -1{
				<-time.Tick(time.Second * 3)
				continue
			}
			fmt.Println(err)
			continue
		}

		self.Scheduler.Submit(*job)

	}
}

func (self *Queue) Pop () (*mtypes.Job,error) {
	job,err := self.Conn.Pop(self.Conn.GetQueueName())
	if err != nil {
		return job,err
	}
	self.Conn.Push(self.Conn.GetQueueExcutingName(),job)
	return job,nil
}

func (self *Queue) Later (second int64,job *mtypes.Job) error {
	return self.Conn.Later(self.Conn.GetQueueName(),second,job)
}

func (self *Queue) Push(job *mtypes.Job) error {
	return self.Conn.Push(self.Conn.GetQueueName(),job)
}

func (self *Queue) Close () error {
	return self.Conn.Close()
}

func createWorker (in chan mtypes.Job,queue *Queue) {
	go func() {
		for {
			job := <- in
			task,has := queue.Handles[job.Handler]
			if has {
				task.Fire(job)
			}
		}
	}()
}