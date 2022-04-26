package valueobj

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

func (r RoleId) Equal(other RoleId) bool {
	return r == other
}
