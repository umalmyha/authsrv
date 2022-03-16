package valueobj

import (
	"errors"
	"fmt"
	"strings"

	"github.com/umalmyha/authsrv/pkg/helpers"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hash string
}

func GeneratePassword(password string, cfg PasswordConfig) (Password, error) {
	var p Password

	if password == "" {
		return p, errors.New("can't be empty")
	}

	if strings.Contains(password, " ") {
		return p, errors.New("spaces are not allowed")
	}

	if cfg.max != 0 && len(password) > cfg.max {
		return p, fmt.Errorf("maximum length is %d characters", cfg.max)
	}

	if len(password) < cfg.min {
		return p, fmt.Errorf("minimum length is %d characters", cfg.min)
	}

	if cfg.hasDigit && !helpers.HasDigit(password) {
		return p, errors.New("must contain at least one digit")
	}

	if cfg.hasUppercase && !helpers.HasUppercase(password) {
		return p, errors.New("must contain at least one uppercase character")
	}

	hash, err := GenerateHash([]byte(password))
	if err != nil {
		return p, err
	}
	p.hash = hash

	return p, nil
}

func GenerateHash(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func PasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

func (p Password) Hash() string {
	return p.hash
}

type PasswordConfig struct {
	min          int
	max          int
	hasDigit     bool
	hasUppercase bool
}

func NewPasswordConfig(min int, max int, hasDigit bool, hasUppercase bool) (PasswordConfig, error) {
	var cfg PasswordConfig

	if max < 0 || min < 0 {
		return cfg, errors.New("minimum and maximum length can't be negative")
	}

	if max != 0 && min > max {
		return cfg, errors.New("minimum length must be less than maximum length")
	}

	return PasswordConfig{
		min:          min,
		max:          max,
		hasDigit:     hasDigit,
		hasUppercase: hasUppercase,
	}, nil
}

func (cfg PasswordConfig) Min() int {
	return cfg.min
}

func (cfg PasswordConfig) Max() int {
	return cfg.max
}

func (cfg PasswordConfig) HasDigit() bool {
	return cfg.hasDigit
}

func (cfg PasswordConfig) HasUppercase() bool {
	return cfg.hasUppercase
}
