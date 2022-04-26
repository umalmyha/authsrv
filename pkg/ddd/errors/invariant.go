package errors

import "encoding/json"

type invariantViolation struct {
	target   string
	msg      string
	severity violationSeverity
}

func NewInvariantViolation(target string, msg string, severity violationSeverity) invariantViolation {
	return invariantViolation{
		target:   target,
		msg:      msg,
		severity: severity,
	}
}

func (iv *invariantViolation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Target   string            `json:"target"`
		Msg      string            `json:"message"`
		Severity violationSeverity `json:"severity"`
	}{
		Target:   iv.target,
		Msg:      iv.msg,
		Severity: iv.severity,
	})
}
