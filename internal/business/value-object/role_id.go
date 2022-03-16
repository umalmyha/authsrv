package valueobj

import (
	"errors"

	"github.com/google/uuid"
)

type RoleId string

func NewRoleId(id string) (RoleId, error) {
	roleId := RoleId(id)

	if _, err := uuid.Parse(id); err != nil {
		return roleId, errors.New("role id must have UUID format")
	}

	return roleId, nil
}

func (r RoleId) String() string {
	return string(r)
}

func (r RoleId) IsTheSameAs(other RoleId) bool {
	return r == other
}