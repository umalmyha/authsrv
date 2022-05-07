package refresh

import (
	"time"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/errors"
)

func NewRefreshToken(fgrprint string, issuedAt time.Time, cfg valueobj.RefreshTokenConfig) (*RefreshToken, error) {
	validation := errors.NewValidation()

	if fgrprint == "" {
		validation.Add(errors.NewBusinessErr("fingerprint is mandatory", "fingerprint", errors.ViolationSeverityErr, errors.CodeValidationFailed))
	}

	if issuedAt.IsZero() {
		validation.Add(errors.NewBusinessErr("issue date can't be initial", "issuedAt", errors.ViolationSeverityErr, errors.CodeValidationFailed))
	}

	if validation.HasError() {
		return nil, pkgerrors.Wrap(validation.RaiseValidationErr(errors.ViolationSeverityErr), "validation failed for refresh token creation")
	}

	expiresAt := issuedAt.Add(cfg.TimeToLive())

	return &RefreshToken{
		id:          uuid.NewString(),
		fingerprint: fgrprint,
		issuedAt:    issuedAt,
		expiresAt:   expiresAt,
	}, nil
}
