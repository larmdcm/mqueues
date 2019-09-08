package mtypes

type Config struct {
	Name string
	Host string
	Port int
	Handlers map[string]Handler
	Connector Connector
	Scheduler Scheduler
	WorkerCount int
}

type Connector interface {
	Delete(queueName string,v interface{}) error
	Pop(queueName string)(*Job,error)
	Later(queueName string,second int64,job *Job) error
	Push(queueName string,job *Job) error
	Connect(config Config) error
	Close() error
	GetQueueName() string
	GetQueueExcutingName() string
	SetQueueName(name string)
}

type Scheduler interface {
	Run()
	Submit(Job)
	WorkeChan() chan Job
}

type Handler interface {
	Fire(job Job)
}