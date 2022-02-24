package valueobj

type NilEmail struct {
	email string
}

func NewNilEmail(s string) (NilEmail, error) {
	ne := NilEmail{email: s}
	if s == "" {
		return ne, nil
	}

	email, err := NewEmail(s)
	if err != nil {
		return ne, err
	}
	ne.email = email.String()

	return ne, nil
}

func NewNilEmailFromPtr(s *string) (NilEmail, error) {
	if s == nil {
		return NewNilEmail("")
	}
	return NewNilEmail(*s)
}

func (ns NilEmail) String() string {
	return ns.email
}

func (ns NilEmail) Ptr() *string {
	s := ns.email
	if s == "" {
		return nil
	}
	return &s
}
