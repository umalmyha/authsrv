package web

type Severity string

const SeverityError Severity = "E"
const SeverityWarning Severity = "W"
const SeveritySuccess Severity = "S"
const SeverityInfo Severity = "I"

type ValidationMessage struct {
	Message  string   `json:"message"`
	Target   string   `json:"target"`
	Severity Severity `json:"severity"`
}

func NewValidationMessage(msg string, target string, severity Severity) ValidationMessage {
	return ValidationMessage{
		Message:  msg,
		Target:   target,
		Severity: severity,
	}
}

type ValidationResult struct {
	valid    bool
	messages []ValidationMessage
}

func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		valid:    true,
		messages: make([]ValidationMessage, 0),
	}
}

func (v *ValidationResult) Add(m ValidationMessage) {
	if m.Severity == SeverityError {
		v.valid = false
	}
	v.messages = append(v.messages, m)
}

func (v *ValidationResult) Valid() bool {
	return v.valid
}

func (v *ValidationResult) Messages() []ValidationMessage {
	return v.messages
}
