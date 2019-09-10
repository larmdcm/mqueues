package handler

import (
	"mqueues"
	"encoding/json"
	"mtypes"
	"time"
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
)

type HttpJsonData struct {

}

type JobHttpHandleConfig struct {
	Url string `json:"url"`
}

type HttpHandler struct {

}

func (self *HttpHandler) Fire(job *mqueues.Job) {
	jobConfig := &JobHttpHandleConfig{}
	json.Unmarshal([]byte(job.JobData.Config),jobConfig)

	jobResult := &mtypes.JobHandleResult{}
	jobResult.Id = job.JobData.Id
	jobResult.Name = job.JobData.Name

	rawData,_ := job.GetJobDataJson()
	jobResult.JobRaw = string(rawData)
	jobResult.Data = job.JobData.Data
	jobResult.Date = time.Now().Unix()

	response,err := self.SendNotify(jobConfig,jobResult)

	if err != nil{
		log.Printf("http handle notify error raw job [%s] request url [%s] errmsg:%s",string(rawData),jobConfig.Url,err.Error())
		return
	}
	log.Printf("http handle notify success response:%s",response)
}

func (self *HttpHandler) SendNotify (config *JobHttpHandleConfig,jobData *mtypes.JobHandleResult) (string,error) {
	jsonData,err := json.Marshal(jobData)
	if err != nil {
		return "",err
	}
	request,err := http.NewRequest(http.MethodPost,config.Url,bytes.NewReader(jsonData))

	if err != nil {
		return "",err
	}

	request.Header.Set("Content-Type","application/json")

	client := &http.Client{}

	response,err := client.Do(request)

	if err != nil {
		return "",err
	}

	defer response.Body.Close()

	body,err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "",err
	}

	return string(body),nil
}