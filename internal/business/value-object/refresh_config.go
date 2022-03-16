package valueobj

import (
	"errors"
	"time"
)

type RefreshTokenConfig struct {
	maxCount int
	ttl      time.Duration
	cookie   string
}

func NewRefreshTokenConfig(ttl time.Duration, count int, cookie string) (RefreshTokenConfig, error) {
	var cfg RefreshTokenConfig

	if ttl == 0 {
		return cfg, errors.New("refresh token ttl must be provided")
	}
	cfg.ttl = ttl

	if count <= 0 {
		return cfg, errors.New("max count can't be zero or negative number")
	}
	cfg.maxCount = count

	if cookie == "" {
		return cfg, errors.New("refresh token cookie name must be provided")
	}
	cfg.cookie = cookie

	return cfg, nil
}

func (cfg RefreshTokenConfig) TimeToLive() time.Duration {
	return cfg.ttl
}

func (cfg RefreshTokenConfig) MaxTokensCount() int {
	return cfg.maxCount
}

func (cfg RefreshTokenConfig) CookieName() string {
	return cfg.cookie
}
