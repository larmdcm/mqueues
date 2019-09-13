package lib

import "os"

func GetBasePath () string {
	path,err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path = path + "/"

	return path + "mweb"
	// return "E:/GoLang/mqueue-tasks/src/" + "mweb"
}

func GetPublicPath () string {
	return GetBasePath() + "/public"
}

func GetViewPath () string {
	return GetBasePath() + "/views"
}