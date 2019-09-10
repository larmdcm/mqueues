package mweb

import (
	"net/http"
	"log"
	"mweb/controller"
	"strings"
	"os"
	"io/ioutil"
	"mweb/lib"
	"mqueues"
)

func staticResource(w http.ResponseWriter, r *http.Request) {
	defer func () {
		if r := recover();r != nil {
			log.Printf("Panic:%v\n",r)
			http.Error(w,http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		}
	}()
	path := r.URL.Path

	if path == "/" {
		siteController  := &controller.SiteController{}
		siteController.Index(w,r)
		return
	}

	request_type := path[strings.LastIndex(path, "."):]
	switch request_type {
	case ".css":
		w.Header().Set("content-type", "text/css")
	case ".js":
		w.Header().Set("content-type", "text/javascript")
	default:
	}
	path = lib.GetPublicPath() + path

	_,err := os.Stat(path)
	if err != nil {
		log.Printf("%s\n",err.Error())
		http.Error(w,http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	if os.IsNotExist(err) {
		log.Printf("%s\n",err.Error())
		http.Error(w,http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}

	fin, err := os.Open(path)

	defer fin.Close()

	if err != nil {
		log.Printf("static resource:%s\n", err.Error())
		http.Error(w,http.StatusText(http.StatusNotFound),http.StatusNotFound)
		return
	}
	fd, _ := ioutil.ReadAll(fin)
	w.Write(fd)
}

func Run (queue *mqueues.Queue)  {
	lib.SetQueue(queue)

	queueController := &controller.QueueController{}
	siteController  := &controller.SiteController{}
	http.HandleFunc("/",staticResource)
	http.HandleFunc("/index",siteController.Index)

	http.HandleFunc("/queue/create",queueController.Create)
	http.HandleFunc("/queue/get",queueController.Get)
	err := http.ListenAndServe(":9099",nil)

	if err != nil {
		log.Fatal(err)
	}
}