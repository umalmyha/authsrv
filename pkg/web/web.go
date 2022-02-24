package web

import (
	"net/http"
)

type httpHandlerFuncWithErr func(http.ResponseWriter, *http.Request) error
type httpErrorHandlerFunc func(http.ResponseWriter, *http.Request, error)

func WithErrorHandler(fn httpHandlerFuncWithErr, errFn httpErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			errFn(w, r, err)
		}
	}
}

func WithDefaultErrorHandler(fn httpHandlerFuncWithErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			defaultErrorHandler(w, err)
		}
	}
}

func defaultErrorHandler(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *RequestErr:
		RespondJson(w, http.StatusBadRequest, err)
	case *NotFoundErr:
		RespondStatus(w, http.StatusNotFound)
	default:
		RespondStatus(w, http.StatusInternalServerError)
	}
}
