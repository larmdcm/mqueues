package mqueues

import "mtypes"

type Config struct {
	Name string
	Handlers map[string]Handler
	Connector Connector
	Scheduler Scheduler
	WorkerCount int
	ConnectConfig mtypes.ConnectConfig
}

type Connector interface {
	Delete(queueName string,v interface{}) error
	Pop(queueName string)(*mtypes.Job,error)
	All(queueName string)(*[]mtypes.Job,error)
	Later(queueName string,delay int64,job *mtypes.Job) error
	Push(queueName string,job *mtypes.Job) error
	Connect(config mtypes.ConnectConfig) error
	Close() error
	GetQueueName() string
	GetQueueExcutingName() string
	SetQueueName(name string)
}

type Scheduler interface {
	Run()
	Submit(mtypes.Job)
	WorkeChan() chan mtypes.Job
}

type Handler interface {
	Fire(job *Job)
}