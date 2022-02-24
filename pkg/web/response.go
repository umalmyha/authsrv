package web

import (
	"encoding/json"
	"net/http"
)

type ResposeWriterFn func(w http.ResponseWriter) error

func RespondJson(w http.ResponseWriter, statusCode int, data interface{}) error {
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(response); err != nil {
		return err
	}

	return nil
}

func RespondStatus(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
