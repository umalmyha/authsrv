package errors

import (
	"fmt"
	"strings"
)

type BusinessError struct {
	errs []*InvariantViolationError
}

func NewBusinessError(errs ...*InvariantViolationError) *BusinessError {
	var e []*InvariantViolationError
	e = append(e, errs...)
	return &BusinessError{
		errs: e,
	}
}

func (e *BusinessError) Error() string {
	var str strings.Builder
	for _, err := range e.errs {
		str.WriteString(fmt.Sprintf("%s. ", err.Error()))
	}
	return str.String()
}
