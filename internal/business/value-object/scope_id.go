package valueobj

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ScopeId string

func NewScopeId(id string) (ScopeId, error) {
	scopeId := ScopeId(id)

	if _, err := uuid.Parse(id); err != nil {
		return scopeId, errors.New("scope id must have UUID format")
	}

	return scopeId, nil
}

func (s ScopeId) String() string {
	return string(s)
}

func (s ScopeId) Equal(other ScopeId) bool {
	return s == other
}
