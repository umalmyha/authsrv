package refresh

import (
	"errors"
	"time"

	"github.com/google/uuid"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	invariant "github.com/umalmyha/authsrv/internal/errors"
)

func NewRefreshToken(fgrprint string, issuedAt time.Time, cfg valueobj.RefreshTokenConfig) (*RefreshToken, error) {
	validation := invariant.NewValidationResult()

	if fgrprint == "" {
		validation.Add(invariant.NewInvariantViolationError("fingerprint is mandatory", "fingerprint"))
		return nil, validation.Error()
	}

	if issuedAt.IsZero() {
		return nil, errors.New("issue date can't be initial")
	}

	expiresAt := issuedAt.Add(cfg.TimeToLive())

	return &RefreshToken{
		id:          uuid.NewString(),
		fingerprint: fgrprint,
		issuedAt:    issuedAt,
		expiresAt:   expiresAt,
	}, nil
}
