package web

import (
	"encoding/json"
	"net/http"
)

type ResposeWriterFn func(w http.ResponseWriter) error

func RespondJson(w http.ResponseWriter, statusCode int, data any) error {
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

func SetHeader(w http.ResponseWriter, header string, value string) {
	w.Header().Set(header, value)
}

func SetCookie(w http.ResponseWriter, cookie *http.Cookie) {
	http.SetCookie(w, cookie)
}

func DeleteCookie(r *http.Request, w http.ResponseWriter, name string) {
	cookie, err := GetCookie(r, name)
	if err != nil {
		return
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
