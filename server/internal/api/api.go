package api

import (
	"encoding/json"
	"log/slog"
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

func ReturnError(log *slog.Logger, w http.ResponseWriter, statusCode int, err error) {
	log.Error("error", slog.Any("error", err))
	w.WriteHeader(statusCode)
	ReturnResponse(w, "error", nil, err.Error())
}
