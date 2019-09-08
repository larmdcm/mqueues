package connection

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"encoding/json"
	"time"
	"errors"
	"mtypes"
)

const (
	QUEUE_NAME 			 = "mqueues:"
	QUEUE_EXECUTING_NAME = ":executeing"
	JOB_NOT_EXISTS 		 = -1
)

type Redis struct {
	Conn redis.Conn
	Name string
}

func (self *Redis) Delete (queueName string,v interface{}) error {
	_,err := self.Conn.Do("ZREM",queueName,v)
	if err != nil {
		return err
	}
	return nil
}

func (self *Redis) Pop (queueName string) (*mtypes.Job,error){
	values,err := redis.Values(self.Conn.Do("ZRANGEBYSCORE",queueName,0,time.Now().Unix()))
	if err != nil {
		return &mtypes.Job{},err
	}
	if len(values) <= 0 {
		return &mtypes.Job{},errors.New(strconv.Itoa(JOB_NOT_EXISTS))
	}
	v := values[0].([]byte)
	err = self.Delete(queueName,v)
	if err != nil {
		return &mtypes.Job{},err
	}
	job := &mtypes.Job{}
	json.Unmarshal(v,job)
	return job,nil
}

func (self *Redis) Later (queueName string,second int64,job *mtypes.Job) (err error) {
	buf,err := json.Marshal(job)
	if err != nil {
		return
	}
	_,err = self.Conn.Do("ZADD",queueName,time.Now().Unix() + second,string(buf))
	return
}

func (self *Redis) Push(queueName string,job *mtypes.Job) error {
	return self.Later(queueName,0,job)
}

func (self *Redis) Connect (config mtypes.Config) error {
	address := config.Host + ":" + strconv.Itoa(config.Port)
	c,err := redis.Dial("tcp",address)

	if err != nil {
		return err
	}
	self.Conn = c
	return nil
}

func (self *Redis) Close () error {
	return self.Conn.Close()
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