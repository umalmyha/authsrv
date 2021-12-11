package validate

import (
	"net/mail"

	"github.com/google/uuid"
)

func Email(email string) (bool, error) {
	_, err := mail.ParseAddress(email)
	return err == nil, err
}

func UUID(s string) (bool, error) {
	_, err := uuid.Parse(s)
	return err == nil, err
}
