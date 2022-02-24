package errors

import (
	"fmt"
)

type InvariantViolationError struct {
	target  string
	message string
}

func NewInvariantViolationError(message string, target string) *InvariantViolationError {
	return &InvariantViolationError{
		target:  target,
		message: message,
	}
}

func (e *InvariantViolationError) Error() string {
	return fmt.Sprintf("[%s] %s", e.target, e.message)
}
