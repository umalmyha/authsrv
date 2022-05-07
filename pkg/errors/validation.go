package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Validation struct {
	errors []*BusinessErr
}

func NewValidation() *Validation {
	return &Validation{
		errors: make([]*BusinessErr, 0),
	}
}

func (v *Validation) Add(errors ...*BusinessErr) {
	v.errors = append(v.errors, errors...)
}

func (v *Validation) Severity() violationSeverity {
	severity := ViolationSeverityInfo
	for _, err := range v.errors {
		if err.severity < severity {
			severity = err.severity
		}
	}
	return severity
}

func (v *Validation) HasSeverity(s violationSeverity) bool {
	for _, err := range v.errors {
		if err.severity == s {
			return true
		}
	}
	return false
}

func (v *Validation) HasError() bool {
	return v.HasSeverity(ViolationSeverityErr)
}

func (v *Validation) HasWarning() bool {
	return v.HasSeverity(ViolationSeverityWarn)
}

func (v *Validation) HasInfo() bool {
	return v.HasSeverity(ViolationSeverityInfo)
}

func (v *Validation) RaiseValidationErr(raiseSeverity violationSeverity) *ValidationErr {
	for _, err := range v.errors {
		if err.severity <= raiseSeverity {
			return &ValidationErr{
				severity: v.Severity(),
				errors:   v.errors,
			}
		}
	}
	return nil
}

type ValidationErr struct {
	severity violationSeverity
	errors   []*BusinessErr
}

func (e *ValidationErr) Error() string {
	var sb strings.Builder
	for _, err := range e.errors {
		sb.WriteString(fmt.Sprintf("%s\n", err.msg))
	}
	return sb.String()
}

func (e *ValidationErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Severity string         `json:"severity"`
		Messages []*BusinessErr `json:"messages"`
	}{
		Severity: e.severity.string(),
		Messages: e.errors,
	})
}
