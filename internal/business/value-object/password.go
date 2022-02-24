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

type PasswordConfig struct {
	Min          int
	Max          int
	HasDigit     bool
	HasUppercase bool
}

func GeneratePassword(password string, cfg PasswordConfig) (Password, error) {
	p := Password{}

	if password == "" {
		return p, errors.New("can't be empty")
	}

	if strings.Contains(password, " ") {
		return p, errors.New("spaces are not allowed")
	}

	if cfg.Min > cfg.Max {
		return p, errors.New("minimum length must be less than maximum length")
	}

	if len(password) < cfg.Min {
		return p, fmt.Errorf("minimum length is %d characters", cfg.Min)
	}

	if len(password) > cfg.Max {
		return p, fmt.Errorf("maximum length is %d characters", cfg.Max)
	}

	if cfg.HasDigit && !helpers.HasDigit(password) {
		return p, errors.New("must contain at least one digit")
	}

	if cfg.HasUppercase && !helpers.HasUppercase(password) {
		return p, errors.New("must contain at least one uppercase character")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	p.hash = string(hash)

	return p, nil
}

func PasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

func (p Password) Hash() string {
	return p.hash
}
