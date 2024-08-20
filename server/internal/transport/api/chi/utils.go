package chiapi

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int      `json:"status"`
	Data   interface{} `json:"data"`
	Error  error      `json:"error,omitempty"`
}

func JSONresponse(w http.ResponseWriter, statusCode int, body interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	res := Response{
		Status: statusCode,
		Data:   body,
		Error:  err,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
