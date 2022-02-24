package valueobj

import "net/mail"

type Email struct {
	str string
}

func NewEmail(s string) (Email, error) {
	email := Email{str: s}
	if _, err := mail.ParseAddress(s); err != nil {
		return email, err
	}
	return email, nil
}

func (e Email) String() string {
	return e.str
}
