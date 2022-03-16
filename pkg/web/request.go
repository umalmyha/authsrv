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

func GetHeader(r *http.Request, header string) string {
	return r.Header.Get(header)
}

func GetCookie(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(name)
}

func GetCookieValue(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
