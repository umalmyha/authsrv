package scope

import (
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
)

type Scope struct {
	id          string
	name        valueobj.SolidString
	description valueobj.NilString
}

func (s *Scope) ChangeDescription(descr string) {
	s.description = valueobj.NewNilString(descr)
}

func (s *Scope) Dto() ScopeDto {
	return ScopeDto{
		Id:          s.id,
		Name:        s.name.String(),
		Description: s.description.Ptr(),
	}
}
