package errors

import (
	"encoding/json"
)

type violationSeverity int

const (
	ViolationSeverityErr violationSeverity = iota + 1
	ViolationSeverityWarn
	ViolationSeverityInfo
)

func (s violationSeverity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string())
}

func (s violationSeverity) string() string {
	switch s {
	case ViolationSeverityErr:
		return "error"
	case ViolationSeverityWarn:
		return "warning"
	default:
		return "info"
	}
}
