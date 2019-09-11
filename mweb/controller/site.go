package controller

import (
	"net/http"
	"mweb/lib"
)

type SiteController struct {
	lib.Controller
}

func (self *SiteController) Index (writer http.ResponseWriter, request *http.Request) {
	self.Render(writer,"index",nil)
}

func (self *SiteController) Create (writer http.ResponseWriter, request *http.Request) {
	self.Render(writer,"create",nil)
}