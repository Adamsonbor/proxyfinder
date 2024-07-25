package api

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}
