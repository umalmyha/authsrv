package refresh

import "time"

type RefreshTokenDto struct {
	Id          string
	Fingerprint string
	CreatedAt   time.Time
}
