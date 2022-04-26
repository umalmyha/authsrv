package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Validation struct {
	violations []invariantViolation
}

func NewValidation() *Validation {
	return &Validation{
		violations: make([]invariantViolation, 0),
	}
}

func (v *Validation) AddViolation(violation invariantViolation) {
	v.violations = append(v.violations, violation)
}

func (v *Validation) AddViolations(violations ...invariantViolation) {
	v.violations = append(v.violations, violations...)
}

func (v *Validation) Severity() violationSeverity {
	severity := ViolationSeverityInfo
	for _, violation := range v.violations {
		if violation.severity < severity {
			severity = violation.severity
		}
	}
	return severity
}

func (v *Validation) HasSeverity(s violationSeverity) bool {
	for _, violation := range v.violations {
		if violation.severity == s {
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

func (v *Validation) Err() *ValidationErr {
	if len(v.violations) == 0 {
		return nil
	}

	return &ValidationErr{
		severity:   v.Severity(),
		violations: v.violations,
	}
}

type ValidationErr struct {
	severity   violationSeverity
	violations []invariantViolation
}

func (e *ValidationErr) Error() string {
	var sb strings.Builder
	lastIndex := len(e.violations) - 1

	for i, violation := range e.violations {
		sb.WriteString(fmt.Sprintf("%s.", violation.msg))
		if i != lastIndex {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}

func (e *ValidationErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Severity string               `json:"severity"`
		Messages []invariantViolation `json:"messages"`
	}{
		Severity: e.severity.string(),
		Messages: e.violations,
	})
}
