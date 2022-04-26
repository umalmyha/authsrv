package refresh

import (
	"time"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/ddd/errors"
)

func NewRefreshToken(fgrprint string, issuedAt time.Time, cfg valueobj.RefreshTokenConfig) (*RefreshToken, error) {
	validation := errors.NewValidation()

	if fgrprint == "" {
		validation.AddViolation(errors.NewInvariantViolation("fingerprint is mandatory", "fingerprint", errors.ViolationSeverityErr))
	}

	if issuedAt.IsZero() {
		validation.AddViolation(errors.NewInvariantViolation("issue date can't be initial", "issuedAt", errors.ViolationSeverityErr))
	}

	if validation.HasError() {
		return nil, pkgerrors.Wrap(validation.Err(), "validation failed for refresh token creation")
	}

	expiresAt := issuedAt.Add(cfg.TimeToLive())

	return &RefreshToken{
		id:          uuid.NewString(),
		fingerprint: fgrprint,
		issuedAt:    issuedAt,
		expiresAt:   expiresAt,
	}, nil
}
