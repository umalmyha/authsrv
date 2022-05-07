package errors

import "encoding/json"

const CodeValidationFailed = "VALID_FAILED"

type BusinessErr struct {
	target   string
	msg      string
	code     string
	severity violationSeverity
}

func NewBusinessErr(target string, msg string, severity violationSeverity, code string) *BusinessErr {
	return &BusinessErr{
		target:   target,
		msg:      msg,
		severity: severity,
		code:     code,
	}
}

func (e *BusinessErr) Error() string {
	return e.msg
}

func (e *BusinessErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Target   string            `json:"target"`
		Msg      string            `json:"message"`
		Severity violationSeverity `json:"severity"`
		Code     string            `json:"code"`
	}{
		Target:   e.target,
		Msg:      e.msg,
		Severity: e.severity,
		Code:     e.code,
	})
}
