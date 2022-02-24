package refresh

import (
	"time"

	"github.com/google/uuid"
	"github.com/umalmyha/authsrv/internal/errors"
)

func NewRefreshToken(fprint string, createdAt time.Time) (*RefreshToken, error) {
	validation := errors.NewValidationResult()

	if fprint == "" {
		validation.Add(errors.NewInvariantViolationError("fingerprint can't be initial", "fingerprint"))
		return nil, validation.Error()
	}

	return &RefreshToken{
		id:          uuid.NewString(),
		fingerprint: fprint,
		createdAt:   createdAt,
	}, nil
}
