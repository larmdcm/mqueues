package mqueues

import (
	"mqueues/connection"
	"testing"
	"time"
	"fmt"
	"log"
	"mtypes"
	"mqueues/scheduler"
)

var mqueue *Queue
var err error

func init ()  {
	config := mtypes.Config{
		Name: "test",
		Host: "127.0.0.1",
		Port: 6379,
		Connector: &connection.Redis{},
		Scheduler: &scheduler.Dispatch{},
		WorkerCount: 50,
	}
	mqueue,err = New(config)
}
func TestQueue_Insert(t *testing.T) {
	job1 := &mtypes.Job{
		time.Now().Unix(),"Push","GoHandle","PushTest","Config",
	}
	job2 := &mtypes.Job{
		time.Now().Unix(),"Later100","HttpHandle","LaterTest","Config",
	}
	err = mqueue.Push(job1)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = mqueue.Later(30,job2)
	fmt.Println(time.Now().Unix())
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestQueue_Pop(t *testing.T) {
	fmt.Println(mqueue.Pop())
}