package connection

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"encoding/json"
	"time"
	"errors"
	"mtypes"
	"mqueues/connection/redispool"
)

const (
	QUEUE_NAME 			 = "mqueues:"
	QUEUE_EXECUTING_NAME = ":executeing"
	JOB_NOT_EXISTS 		 = -1
)

type Redis struct {
	Pool *redis.Pool
	Name string

	Host string
	Port int
	PassWord string
}

func (self *Redis) Delete (queueName string,v interface{}) error {
	redisConn := self.Pool.Get()
	defer redisConn.Close()

	_,err := redisConn.Do("ZREM",queueName,v.([]byte))
	if err != nil {
		return err
	}
	return nil
}

func (self *Redis) Pop (queueName string) (*mtypes.Job,error) {
	redisConn := self.Pool.Get()
	defer redisConn.Close()

	values,err := redis.Values(redisConn.Do("ZRANGEBYSCORE",queueName,0,time.Now().Unix()))
	if err != nil {
		return &mtypes.Job{},err
	}
	if len(values) <= 0 {
		return &mtypes.Job{},errors.New(strconv.Itoa(JOB_NOT_EXISTS))
	}
	err = self.Delete(queueName,values[0])
	v := values[0].([]byte)
	if err != nil {
		return &mtypes.Job{},err
	}
	job := &mtypes.Job{}
	json.Unmarshal(v,job)
	return job,nil
}

func (self *Redis) All (queueName string) (*[]mtypes.Job,error) {
	redisConn := self.Pool.Get()
	defer redisConn.Close()
	jobs := []mtypes.Job{}
	values,err := redis.Values(redisConn.Do("ZRANGE",queueName,0,-1))
	if err != nil {
		return &jobs,err
	}

	if len(values) <= 0 {
		return &jobs,errors.New(strconv.Itoa(JOB_NOT_EXISTS))
	}

	for _,v := range values {
		job := &mtypes.Job{}
		err := json.Unmarshal(v.([]byte),job)
		if err != nil {
			continue
		}
		jobs = append(jobs,*job)
	}
	return &jobs,nil
}

func (self *Redis) Later (queueName string,delay int64,job *mtypes.Job) (err error) {
	redisConn := self.Pool.Get()
	defer redisConn.Close()

	buf,err := json.Marshal(job)
	if err != nil {
		return
	}
	_,err = redisConn.Do("ZADD",queueName,time.Now().Unix() + delay,string(buf))
	return
}

func (self *Redis) Push(queueName string,job *mtypes.Job) error {
	return self.Later(queueName,0,job)
}

func (self *Redis) Connect (config mtypes.ConnectConfig) error {
	self.Host 	  = config.Host
	self.Port 	  = config.Port
	self.PassWord = config.PassWord

	self.Pool = self.newPool()
	return nil
}

func (self *Redis) Close () error {
	return self.Pool.Close()
}

func (self *Redis) GetQueueName () string {
	return QUEUE_NAME + self.Name
}

func (self *Redis) GetQueueExcutingName () string {
	return QUEUE_NAME + self.Name + QUEUE_EXECUTING_NAME
}

func (self *Redis) SetQueueName (name string) {
	self.Name = name
}

func (self *Redis) newPool () *redis.Pool {
	return redispool.NewPool(self.Host,strconv.Itoa(self.Port),self.PassWord)
}