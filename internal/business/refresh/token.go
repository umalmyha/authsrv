package refresh

import (
	"time"

	"github.com/umalmyha/authsrv/pkg/errors"
)

var RefreshTokenExpiredErr = errors.NewBusinessErr(
	"refreshToken",
	"refresh token already expired",
	errors.ViolationSeverityErr,
	"RFR_TOKEN_EXPIRED",
)

type RefreshToken struct {
	id          string
	fingerprint string
	issuedAt    time.Time
	expiresAt   time.Time
}

func (rt *RefreshToken) Id() string {
	return rt.id
}

func (rt *RefreshToken) Fingerprint() string {
	return rt.fingerprint
}

func (rt *RefreshToken) IssuedAt() time.Time {
	return rt.issuedAt
}

func (rt *RefreshToken) ExpiresAt() time.Time {
	return rt.expiresAt
}

func (rt *RefreshToken) UnixExpiresIn() int {
	return int(rt.expiresAt.Unix() - rt.issuedAt.Unix())
}

func (rt *RefreshToken) VerifyNotExpired(now time.Time) error {
	if rt.ExpiresAt().Before(now) {
		return RefreshTokenExpiredErr
	}
	return nil
}
