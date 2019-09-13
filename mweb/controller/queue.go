package controller

import (
	"net/http"
	"mtypes"
	"mweb/lib"
	"encoding/json"
	"log"
	"strconv"
	"mqueues"
)

type QueueController struct {
	lib.Controller
}

type QueueJobListResult struct {
	QueueJobs string `json:"queue_jobs"`
	QueueJobExcuteings string `json:"queue_job_excuteings"`
}

func (self *QueueController) Release (writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type","application/json")

	if request.Method != http.MethodPost {
		self.HttpJsonError(writer,"request method is not a post")
		return
	}

	jobRaw := request.FormValue("job_raw")
	delayStr := request.FormValue("delay")

	delay,err := strconv.ParseInt(delayStr,10,64)
	if err != nil {
		self.HttpJsonError(writer,"delay type is not a int")
		return
	}
	jobData := &mtypes.Job{}
	json.Unmarshal([]byte(jobRaw),jobData)
	job := &mqueues.Job{
		JobData: jobData,
		Queue: lib.GetQueue(),
	}
	err = job.Release(delay)
	if err != nil {
		log.Printf("queue/release jobs release error:%s",err.Error())
		self.Error(writer,"queue job release error","{}")
		return
	}
	self.Correct(writer,"queue job release success","{}")
}
func (self *QueueController) Delete (writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type","application/json")

	if request.Method != http.MethodPost {
		self.HttpJsonError(writer,"request method is not a post")
		return
	}

	queueType := request.FormValue("queue_type")
	jobRaw := request.FormValue("job_raw")
	conn := lib.GetQueue().Conn
	queueName := conn.GetQueueName()
	if queueType != "1" {
		queueName = conn.GetQueueExcutingName()
	}
	err := conn.Delete(queueName,[]byte(jobRaw))
	if err != nil {
		log.Printf("queue/delete jobs delete error:%s",err.Error())
		self.Error(writer,"queue job delete error","{}")
		return
	}
	self.Correct(writer,"queue job delete success","{}")
}

func (self *QueueController) Get (writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type","application/json")

	result := &QueueJobListResult{}

	getJobs := func(queueName string) string {
		jobs,err := lib.GetQueue().Conn.All(queueName)
		if err != nil {
			if err.Error() == "-1" {
				return "[]"
			}
			log.Printf("queue/get jobs get error:%s",err.Error())
			return "[]"
		}
		response,err := json.Marshal(jobs)
		if err != nil {
			log.Printf("get queue %s json format error -> %s\n",queueName,err.Error())
			return "[]"
		}
		return string(response)
	}

	result.QueueJobs = getJobs(lib.GetQueue().Conn.GetQueueName())
	result.QueueJobExcuteings = getJobs(lib.GetQueue().Conn.GetQueueExcutingName())

	response,err := json.Marshal(result)

	if err != nil {
		log.Printf("response json format error -> %s\n",err.Error())
		self.HttpError(writer,"json format error")
		return
	}
	self.Correct(writer,"get success",string(response))

}

func (self *QueueController) Create (writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type","application/json")

	if request.Method != http.MethodPost {
		self.HttpJsonError(writer,"request method is not a post")
		return
	}
	
	id := mqueues.UniqueId()
	name := request.FormValue("name")
	handler := request.FormValue("handler")
	data := request.FormValue("data")
	config := request.FormValue("config")
	delayStr := request.FormValue("delay")

	delay,err := strconv.ParseInt(delayStr,10,64)
	if err != nil {
		self.HttpJsonError(writer,"delay type is not a int")
		return
	}
	job := &mtypes.Job{
		Id: id,
		Name: name,
		Handler: handler,
		Data: data,
		Config: config,
		AttemptsCount: 0,
	}

	err = lib.GetQueue().Later(delay,job)
	var code int64 = 0
	msg := "successful delivery"
	if err != nil {
		msg ="Delivery failure"
		code = 400
	}
	self.ResultJson(writer,code,msg,"{}")
}