package lib

import (
	"fmt"
	"html/template"
)

var tplPath string

func init ()  {
	tplPath = GetViewPath()
}

func Fetch (view string) (*template.Template,error) {
	tpl,err := template.ParseFiles(tplPath + "/" + view + ".html")

	if err != nil {
		fmt.Println("parse file err:", err)
		return &template.Template{},err
	}
	return tpl,err
}
