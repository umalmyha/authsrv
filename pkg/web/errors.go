package web

import (
	"fmt"
	"strings"
)

type RequestErrMessage struct {
	Message string `json:"message"`
	Target  string `json:"target"`
}

func NewRequestErrMessage(msg string, target string) RequestErrMessage {
	return RequestErrMessage{
		Message: msg,
		Target:  target,
	}
}

type RequestErr struct {
	Messages []RequestErrMessage `json:"messages"`
}

func NewRequestErr(messages ...RequestErrMessage) error {
	return &RequestErr{
		Messages: messages,
	}
}

func (e *RequestErr) Error() string {
	var b strings.Builder

	for _, msg := range e.Messages {
		b.WriteString(fmt.Sprintf("%s. ", msg.Message))
	}

	return b.String()
}

type NotFoundErr struct{}

func NewNotFoundError(msg string) error {
	return &NotFoundErr{}
}

func (e *NotFoundErr) Error() string {
	return "Not Found"
}
