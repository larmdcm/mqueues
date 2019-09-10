package main

import (
	"mqueues/connection"
	"mqueues/scheduler"
	"mqueues"
	"mtypes"
	"handler"
	"mweb"
)

var mqueue *mqueues.Queue
var err error
var handlers = make(map[string]mqueues.Handler)

func init ()  {
	handlers["GoHandle"] = &handler.GoHandler{}
	handlers["HttpHandle"] = &handler.HttpHandler{}
	config := mqueues.Config{
		Name: "test",
		ConnectConfig: mtypes.ConnectConfig{
			Host: "127.0.0.1",
			Port: 6379,
			PassWord: "",
		},
		Connector: &connection.Redis{},
		Scheduler: &scheduler.Dispatch{},
		WorkerCount: 100,
		Handlers: handlers,
	}
	mqueue,err = mqueues.New(config)
}

func main() {
	go mqueue.Run()
	mweb.Run(mqueue)
}
