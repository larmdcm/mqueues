package handler

import (
	"mqueues"
	"log"
)

type GoHandler struct {

}

func (self *GoHandler) Fire(job *mqueues.Job) {
	if job.Attempts() > 1 {
		job.Delete()
		return
	}
	log.Printf("GoHandler job run attempts count:%d\n",job.Attempts())
	job.Release(0)
}