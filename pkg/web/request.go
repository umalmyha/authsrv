package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func UrlParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

func JsonReqBody(r *http.Request, to interface{}) error {
	return json.NewDecoder(r.Body).Decode(to)
}
