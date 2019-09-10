package connection

import (
	"testing"
	"log"
	"fmt"
	"time"
	"mtypes"
)

var mredis *Redis

func init () {
	mredis = &Redis{Name: "test"}

	err := mredis.Connect(mtypes.ConnectConfig{
		Host: "127.0.0.1",
		Port: 6379,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestRedis_Insert(t *testing.T) {

	var err error
	defer mredis.Close()

	job1 := &mtypes.Job{
		time.Now().Unix(),"Push","Handler","PushTest","test",0,
	}
	job2 := &mtypes.Job{
		time.Now().Unix(),mredis.GetQueueName(),"HttpHandler","LaterTest","test",0,
	}
	err = mredis.Push("test",job1)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = mredis.Later(mredis.GetQueueName(),100,job2)
	fmt.Println(time.Now().Unix())
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestRedis_Pop(t *testing.T) {
	defer mredis.Close()
	job,err := mredis.Pop(mredis.GetQueueName())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(job)
}