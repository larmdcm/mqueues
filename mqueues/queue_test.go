package mqueues

import (
	"mqueues/connection"
	"testing"
	"time"
	"fmt"
	"mtypes"
	"mqueues/scheduler"
)

var mqueue *Queue
var err error

func init ()  {
	config := Config{
		Name: "test",
		ConnectConfig: mtypes.ConnectConfig{
			"127.0.0.1",6379,"",
		},
		Connector: &connection.Redis{},
		Scheduler: &scheduler.Dispatch{},
		WorkerCount: 50,
	}
	mqueue,err = New(config)
}
func TestQueue_Insert(t *testing.T) {
	job1 := &mtypes.Job{
		"1","Push","GoHandle","PushTest","Config",0,
	}
	job2 := &mtypes.Job{
		"2","Later100","HttpHandle","LaterTest","Config",0,
	}
	err = mqueue.Push(job1)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = mqueue.Later(30,job2)
	fmt.Println(time.Now().Unix())
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestQueue_Pop(t *testing.T) {
	fmt.Println(mqueue.Pop())
}

func TestJob_GetJobDataJson(t *testing.T) {
	job,err := mqueue.Pop()
	if err != nil {
		t.Fatal(err)
	}
	mjob := &Job{
		job,mqueue,
	}
	b,err := mjob.GetJobDataJson()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestJob_Delete(t *testing.T) {
	job,err := mqueue.Pop()
	if err != nil {
		t.Fatal(err)
	}
	mjob := &Job{
		job,mqueue,
	}
	err = mjob.Delete()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("queue job delete at success")
}

func TestJob_Release(t *testing.T) {
	job,err := mqueue.Pop()
	if err != nil {
		t.Fatal(err)
	}
	mjob := &Job{
		job,mqueue,
	}
	err = mjob.Release(0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("queue job release at success job attempts count:%d",mjob.Attempts())
}