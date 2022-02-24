package user

import (
	"container/list"
	"fmt"

	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
)

type roleExistFn func(string) (bool, error)

type User struct {
	id            string
	username      valueobj.SolidString
	email         valueobj.NilEmail
	password      valueobj.Password
	isSuperuser   bool
	firstName     valueobj.NilString
	lastName      valueobj.NilString
	middleName    valueobj.NilString
	assignedRoles *list.List
}

func (u *User) AssignRole(roleId string, existFn roleExistFn) error {
	roleToAssign, err := valueobj.NewAssignedRole(u.id, roleId)
	if err != nil {
		return err
	}

	for elem := u.assignedRoles.Front(); elem != nil; elem = elem.Next() {
		role, _ := elem.Value.(valueobj.AssignedRole)
		if role.IsTheSameAs(roleToAssign) {
			return fmt.Errorf("role with id %s is already assigned", roleId)
		}
	}

	if exist, err := existFn(roleId); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("role with id %s doesn't exist", roleId)
	}

	u.assignedRoles.PushBack(roleToAssign)
	return nil
}

func (u *User) UnassignRole(roleId string) error {
	roleToUnassign, err := valueobj.NewAssignedRole(u.id, roleId)
	if err != nil {
		return err
	}

	var rmElem *list.Element
	for elem := u.assignedRoles.Front(); elem != nil; elem = elem.Next() {
		role, _ := elem.Value.(valueobj.AssignedRole)
		if role.IsTheSameAs(roleToUnassign) {
			rmElem = elem
			break
		}
	}

	if rmElem == nil {
		return fmt.Errorf("role with id %s is not assigned", roleId)
	}

	u.assignedRoles.Remove(rmElem)
	return nil
}

func (u *User) ToDto() UserDto {
	return UserDto{
		Id:          u.id,
		Username:    u.username.String(),
		Email:       u.email.Ptr(),
		Password:    u.password.Hash(),
		IsSuperuser: u.isSuperuser,
		FirstName:   u.firstName.Ptr(),
		LastName:    u.lastName.Ptr(),
		MiddleName:  u.middleName.Ptr(),
	}
}

func (u *User) RolesDto() []AssignedRoleDto {
	dto := make([]AssignedRoleDto, 0)

	for elem := u.assignedRoles.Front(); elem != nil; elem = elem.Next() {
		role, _ := elem.Value.(valueobj.AssignedRole)
		roleDto := AssignedRoleDto{UserId: role.UserId(), RoleId: role.RoleId()}
		dto = append(dto, roleDto)
	}

	return dto
}
