package valueobj

type NilString struct {
	str string
}

func NewNilString(s string) NilString {
	return NilString{
		str: s,
	}
}

func NewNilStringFromPtr(s *string) NilString {
	if s == nil {
		return NewNilString("")
	}
	return NewNilString(*s)
}

func (ns NilString) String() string {
	return ns.str
}

func (ns NilString) Ptr() *string {
	s := ns.str
	if s == "" {
		return nil
	}
	return &s
}
