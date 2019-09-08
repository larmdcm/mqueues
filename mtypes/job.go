package mtypes

type Job struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Handler string `json:"handler"`
	Data string `json:"data"`
	Config string `json:"config"`
}