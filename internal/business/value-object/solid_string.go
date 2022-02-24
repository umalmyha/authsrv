package valueobj

import (
	"errors"
	"regexp"
)

type SolidString struct {
	str string
}

func NewSolidString(s string) (SolidString, error) {
	str := SolidString{str: s}
	if ok, err := regexp.MatchString(`^\S+$`, s); err != nil {
		return str, err
	} else if !ok {
		return str, errors.New("spaces are not allowed")
	}
	return str, nil
}

func (s SolidString) String() string {
	return s.str
}
