package web

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Err      string              `json:"error"`
	Messages []ValidationMessage `json:"messages"`
}

func NewValidationError(messages ...ValidationMessage) error {
	return &ValidationError{
		Err:      "Validation failed",
		Messages: messages,
	}
}

func (e *ValidationError) Error() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s .", e.Err))

	for _, msg := range e.Messages {
		b.WriteString(fmt.Sprintf("%s .", msg.Message))
	}

	return b.String()
}

type NotFoundError struct {
	Err string
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{
		Err: msg,
	}
}

func (e *NotFoundError) Error() string {
	return e.Err
}
