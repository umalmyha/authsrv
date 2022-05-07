package web

import (
	"errors"
	"net/http"

	webErrs "github.com/umalmyha/authsrv/pkg/web/errors"
	"github.com/umalmyha/authsrv/pkg/web/response"
)

type HttpHandlerFn func(http.ResponseWriter, *http.Request) error
type HttpErrorHandlerFn func(http.ResponseWriter, *http.Request, error)

func HttpHandlerFunc(handlerFn HttpHandlerFn) http.HandlerFunc {
	return HttpFuncWithErrHandler(handlerFn, DefaultErrorHandler)
}

func HttpFuncWithErrHandler(handlerFn HttpHandlerFn, errHandlerFn HttpErrorHandlerFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handlerFn(w, r); err != nil {
			errHandlerFn(w, r, err)
		}
	}
}

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	switch e := err.(type) {
	case *webErrs.HttpErrWithBody:
		handleHttpErrWithBody(w, e)
	case *webErrs.HttpErr:
		response.RespondStatus(w, e.Status())
	default:
		response.RespondStatus(w, http.StatusInternalServerError)
	}
}

func handleHttpErrWithBody(w http.ResponseWriter, httpErr *webErrs.HttpErrWithBody) {
	var err error
	switch httpErr.ContentType() {
	case "application/json":
		err = response.RespondJson(w, httpErr.Status(), httpErr.Data())
	default:
		b, ok := httpErr.Data().([]byte)
		if ok {
			err = response.RespondTextPlain(w, httpErr.Status(), b)
		} else {
			err = errors.New("failed to write text/plain content because provided data is not serialized")
		}
	}

	if err != nil {
		panic(err)
	}
}
