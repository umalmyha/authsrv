package valueobj

import (
	"crypto/rsa"
	"time"

	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Jwt struct {
	signed    string
	tokenType string
	expiresAt time.Time
}

func NewJwt(user string, issuedAt time.Time, roles []string, scopes []string, cfg JwtConfig) (Jwt, error) {
	var accessToken Jwt

	if user == "" {
		return accessToken, errors.New("user is mandatory for JWT generation (used as subject)")
	}

	if roles == nil {
		roles = make([]string, 0)
	}

	if scopes == nil {
		scopes = make([]string, 0)
	}

	method := jwt.GetSigningMethod(cfg.algorithm)

	if issuedAt.IsZero() {
		return accessToken, errors.New("issued timestamp is mandatory")
	}

	expiresAt := issuedAt.Add(cfg.ttl)
	accessToken.expiresAt = expiresAt

	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    cfg.issuer,
			Subject:   user,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		SubjRoles:  roles,
		SubjScopes: scopes,
	}

	token := jwt.NewWithClaims(method, claims)
	signed, err := token.SignedString(cfg.privateKey)
	if err != nil {
		return accessToken, err
	}

	accessToken.signed = signed
	accessToken.tokenType = "Bearer"

	return accessToken, nil
}

func (jwt Jwt) String() string {
	return jwt.signed
}

func (jwt Jwt) ExpiresAt() int64 {
	return jwt.expiresAt.Unix()
}

func (jwt Jwt) TokenType() string {
	return jwt.tokenType
}

type JwtClaims struct {
	jwt.RegisteredClaims
	SubjRoles  []string `json:"roles"`
	SubjScopes []string `json:"scopes"`
}

func (c JwtClaims) Username() string {
	return c.Subject
}

func (c JwtClaims) Roles() []string {
	return c.SubjRoles
}

func (c JwtClaims) Scopes() []string {
	return c.SubjScopes
}

type JwtConfig struct {
	algorithm  string
	issuer     string
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	ttl        time.Duration
}

func NewJwtConfig(alg string, issuer string, rsaPrivate *rsa.PrivateKey, rsaPublic *rsa.PublicKey, ttl time.Duration) (JwtConfig, error) {
	var cfg JwtConfig

	if jwt.GetSigningMethod(alg) == nil {
		return cfg, errors.Errorf("%s invalid alogrithm for JWT generation", alg)
	}
	cfg.algorithm = alg

	if rsaPrivate == nil {
		return cfg, errors.New("private key can't be initial")
	}
	cfg.privateKey = rsaPrivate

	if rsaPublic == nil {
		return cfg, errors.New("public key can't be initial")
	}
	cfg.publicKey = rsaPublic

	if issuer == "" {
		return cfg, errors.New("issuer can't be initial")
	}
	cfg.issuer = issuer

	if ttl == 0 {
		return cfg, errors.New("ttl must be provided")
	}
	cfg.ttl = ttl

	return cfg, nil
}

func (cfg JwtConfig) Algorithm() string {
	return cfg.algorithm
}

func (cfg JwtConfig) Issuer() string {
	return cfg.issuer
}

func (cfg JwtConfig) PrivateKey() *rsa.PrivateKey {
	return cfg.privateKey
}

func (cfg JwtConfig) PublicKey() *rsa.PublicKey {
	return cfg.publicKey
}

func (cfg JwtConfig) TimeToLive() time.Duration {
	return cfg.ttl
}
