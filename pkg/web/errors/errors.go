package errors

import (
	"net/http"
)

var (
	HttpNotFoundErr       = NewHttpErr(http.StatusNotFound, "Not Found")
	HttpUnauthorizedErr   = NewHttpErr(http.StatusUnauthorized, "Unauthorized")
	HttpForbiddenErr      = NewHttpErr(http.StatusForbidden, "Forbidden")
	HttpInternalServerErr = NewHttpErr(http.StatusInternalServerError, "Internal Server Error")
)

func HttpBadRequestErr(cType string, body any) error {
	return &HttpErrWithBody{
		HttpErr:     NewHttpErr(http.StatusBadRequest, "Bad Request"),
		contentType: cType,
		body:        body,
	}
}

func HttpBadRequestJsonErr(body any) error {
	return HttpBadRequestErr("application/json", body)
}

type HttpErr struct {
	status  int
	message string
}

func NewHttpErr(status int, msg string) *HttpErr {
	return &HttpErr{
		status:  status,
		message: msg,
	}
}

func (e *HttpErr) Error() string {
	return e.message
}

func (e *HttpErr) Status() int {
	return e.status
}

type HttpErrWithBody struct {
	*HttpErr
	body        any
	contentType string
}

func (e *HttpErrWithBody) Data() any {
	return e.body
}

func (e *HttpErrWithBody) ContentType() string {
	return e.contentType
}
