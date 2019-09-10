package lib

import (
	"net/http"
	"fmt"
	"mweb/model"
	"time"
	"log"
	"encoding/json"
)

type Controller struct {

}

func (self *Controller) HttpError (writer http.ResponseWriter,error string) {
	http.Error(writer,error,http.StatusInternalServerError)
}

func (self *Controller) Render (writer http.ResponseWriter,view string,data interface{}) {
	tpl,err := Fetch(view)
	if err != nil {
		self.HttpError(writer,fmt.Sprintf("view file fetch error"))
		return
	}
	tpl.Execute(writer,data)
}

func (self *Controller) HttpJsonError (writer http.ResponseWriter,msg string) {
	 resultData := model.ResultData{
	 	Code: 500,
	 	Data: "[]",
	 	Msg: msg,
	 	Date: time.Now().Unix(),
	 }
	 response,err := json.Marshal(resultData)
	 if err != nil {
		log.Printf("response json format error -> %s\n",err.Error())
		return
	 }
	 writer.Write(response)
}

func (self *Controller) ResultJson (writer http.ResponseWriter,code int64,msg string,data string) {
	resultData := &model.ResultData{
		Date: time.Now().Unix(),
		Data: data,
		Code: code,
		Msg: msg,
	}
	response,err := json.Marshal(resultData)
	if err != nil {
		log.Printf("response json format error -> %s\n",err.Error())
		self.HttpJsonError(writer,"result data json format error")
		return
	}
	writer.Write(response)
}

func (self *Controller) Correct (writer http.ResponseWriter,msg string,data string) {
	self.ResultJson(writer,0,msg,data)
}

func (self *Controller) Error (writer http.ResponseWriter,msg string,data string) {
	self.ResultJson(writer,400,msg,data)
}