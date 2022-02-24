package valueobj

import "errors"

type AssignedRole struct {
	userId string
	roleId string
}

func NewAssignedRole(userId, roleId string) (AssignedRole, error) {
	role := AssignedRole{
		userId: userId,
		roleId: roleId,
	}

	if userId == "" {
		return role, errors.New("user id must be provided")
	}

	if roleId == "" {
		return role, errors.New("role id must be provided")
	}

	return role, nil
}

func (r AssignedRole) IsTheSameAs(other AssignedRole) bool {
	return r.roleId == other.roleId && r.userId == other.userId
}

func (r AssignedRole) UserId() string {
	return r.userId
}

func (r AssignedRole) RoleId() string {
	return r.roleId
}
