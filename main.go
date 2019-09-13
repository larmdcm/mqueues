package main

import (
	"mqueues/connection"
	"mqueues/scheduler"
	"mqueues"
	"mtypes"
	"handler"
	"os"
	"log"
	"github.com/Unknwon/goconfig"
	"mweb"
)

var mqueue *mqueues.Queue
var err error
var handlers = make(map[string]mqueues.Handler)
var fconfig *goconfig.ConfigFile

func init ()  {
	path ,_ := os.Getwd()
	fconfig,err = goconfig.LoadConfigFile(path + "/mqueues.conf")
	if err != nil {
		log.Fatal("config read error:%s",err.Error())
	}
	queueName,_ := fconfig.GetValue("queue","name")
	redisHost,_ := fconfig.GetValue("redis","host")
	redisPassword,_ := fconfig.GetValue("redis","password")
	redisPort,_ := fconfig.Int("redis","port")

	handlers["GoHandle"] = &handler.GoHandler{}
	handlers["HttpHandle"] = &handler.HttpHandler{}
	config := mqueues.Config{
		Name: queueName,
		ConnectConfig: mtypes.ConnectConfig{
			Host: redisHost,
			Port: redisPort,
			PassWord: redisPassword,
		},
		Connector: &connection.Redis{},
		Scheduler: &scheduler.Dispatch{},
		WorkerCount: 100,
		Handlers: handlers,
	}
	mqueue,err = mqueues.New(config)
}

func main() {
	out := make(chan bool)
	go mqueue.Run()

	webServerIsStart,_ := fconfig.Int("web_server","start")
	if webServerIsStart == 1{
		webHost,_ := fconfig.GetValue("web_server","host")
		webPort,_ := fconfig.GetValue("web_server","port")
		webConig := map[string]string {
			"host": webHost,
			"port": webPort,
		}
		go mweb.Run(mqueue,webConig)
	}
	<-out
}
