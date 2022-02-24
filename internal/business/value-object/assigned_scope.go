package valueobj

import "errors"

type AssignedScope struct {
	scopeId string
	roleId  string
}

func NewAssignedScope(roleId, scopeId string) (AssignedScope, error) {
	scope := AssignedScope{scopeId: scopeId, roleId: roleId}

	if scopeId == "" {
		return scope, errors.New("scope id must be provided")
	}

	if roleId == "" {
		return scope, errors.New("role id must be provided")
	}

	return scope, nil
}

func (s AssignedScope) RoleId() string {
	return s.roleId
}

func (s AssignedScope) ScopeId() string {
	return s.scopeId
}

func (s AssignedScope) IsTheSameAs(other AssignedScope) bool {
	return s.roleId == other.roleId && s.scopeId == other.scopeId
}
