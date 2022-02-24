package refresh

import "time"

type RefreshToken struct {
	id          string
	fingerprint string
	createdAt   time.Time
}

func (token *RefreshToken) ToDto() RefreshTokenDto {
	return RefreshTokenDto{
		Id:          token.id,
		Fingerprint: token.fingerprint,
		CreatedAt:   token.createdAt,
	}
}
