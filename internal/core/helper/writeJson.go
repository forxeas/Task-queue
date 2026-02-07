package helper

import (
	"encoding/json"
	"net/http"
	"time"
)

type JsonError struct {
	Errors string    `json:"errors"`
	Time   time.Time `json:"time"`
}

func NewJsonError(errors error) *JsonError {
	return &JsonError{Errors: errors.Error(), Time: time.Now()}
}

func WriteJson(w http.ResponseWriter, jobs any, httpCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)

	return json.NewEncoder(w).Encode(jobs)
}

func WriteJsonError(w http.ResponseWriter, err error, httpCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)

	return json.NewEncoder(w).Encode(NewJsonError(err))
}
