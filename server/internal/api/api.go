package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

func ReturnResponse(
	w http.ResponseWriter,
	status string,
	data interface{},
	err string) {
	res := Response{
		Status: status,
		Data:   data,
		Error:  err,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
