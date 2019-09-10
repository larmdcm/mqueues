package mtypes

type Job struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Handler string `json:"handler"`
	Data string `json:"data"`
	Config string `json:"config"`
	AttemptsCount int64 `json:"attempts_count"`
}

type JobHandleResult struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Data string `json:"data"`
	JobRaw string `json:"job_raw"`
	Date int64 `json:"date"`
}