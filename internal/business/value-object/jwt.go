package valueobj

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const defaultIssuer = "umalmyha/authsrv"
const defaultTtl = 10 * time.Minute

type Jwt string

func NewJwt(userId string, roles []string, scopes []string, cfg JwtConfig) (Jwt, error) {
	if userId == "" {
		return "", errors.New("user is mandatory for JWT generation (used as issuer)")
	}

	if roles == nil {
		roles = make([]string, 0)
	}

	if scopes == nil {
		scopes = make([]string, 0)
	}

	method := jwt.GetSigningMethod(cfg.algorithm)

	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    cfg.issuer,
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.ttl)),
		},
		Roles:  roles,
		Scopes: scopes,
	}

	token := jwt.NewWithClaims(method, claims)
	signed, err := token.SignedString(cfg.privateKey)
	if err != nil {
		return "", err
	}

	return Jwt(signed), nil
}

type JwtClaims struct {
	jwt.RegisteredClaims
	Roles  []string `json:"roles"`
	Scopes []string `json:"scopes"`
}

type JwtConfig struct {
	algorithm  string
	issuer     string
	privateKey string
	ttl        time.Duration
}

func NewJwtConfig(alg string, pkey string, issuer string, ttl time.Duration) (JwtConfig, error) {
	var cfg JwtConfig

	if jwt.GetSigningMethod(alg) == nil {
		return cfg, errors.New(fmt.Sprintf("%s invalid alogrithm for JWT generation", alg))
	}
	cfg.algorithm = alg

	if pkey == "" {
		return cfg, errors.New("private key can't be initial")
	}
	cfg.privateKey = pkey

	cfg.issuer = issuer
	if cfg.issuer == "" {
		cfg.issuer = defaultIssuer
	}

	cfg.ttl = ttl
	if cfg.ttl == 0 {
		cfg.ttl = defaultTtl
	}

	return cfg, nil
}
