package main

import (
	"mqueues/connection"
	"mtypes"
	"mqueues/scheduler"
	"mqueues"
)

var mqueue *mqueues.Queue
var err error

func init ()  {
	config := mtypes.Config{
		Name: "test",
		Host: "127.0.0.1",
		Port: 6379,
		Connector: &connection.Redis{},
		Scheduler: &scheduler.Dispatch{},
		WorkerCount: 100,
	}
	mqueue,err = mqueues.New(config)
}

func main() {
	mqueue.Run()
}
