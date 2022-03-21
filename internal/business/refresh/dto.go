package refresh

import "time"

type RefreshTokenDto struct {
	Id          string
	Fingerprint string
	UserId      string
	IssuedAt    time.Time
	ExpiresAt   time.Time
}

func (dto RefreshTokenDto) Key() string {
	return dto.Id
}

func (dto RefreshTokenDto) IsPresent() bool {
	return dto.Id != ""
}

func (dto RefreshTokenDto) Equal(other RefreshTokenDto) bool {
	return dto.Id == other.Id &&
		dto.Fingerprint == other.Fingerprint &&
		dto.IssuedAt == other.IssuedAt &&
		dto.ExpiresAt == other.ExpiresAt
}

func (dto RefreshTokenDto) Clone() RefreshTokenDto {
	return dto
}

func (dto RefreshTokenDto) ToRefreshToken() *RefreshToken {
	return &RefreshToken{
		id:          dto.Id,
		fingerprint: dto.Fingerprint,
		issuedAt:    dto.IssuedAt,
		expiresAt:   dto.ExpiresAt,
	}
}
