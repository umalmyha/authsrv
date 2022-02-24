package errors

type validationResult struct {
	errs []*InvariantViolationError
}

func NewValidationResult() *validationResult {
	return &validationResult{
		errs: make([]*InvariantViolationError, 0),
	}
}

func (e *validationResult) Add(err *InvariantViolationError) {
	e.errs = append(e.errs, err)
}

func (e *validationResult) Error() error {
	if e.Failed() {
		return NewBusinessError(e.errs...)
	}
	return nil
}

func (e *validationResult) Failed() bool {
	return len(e.errs) > 0
}
