package role

import (
	"fmt"

	"container/list"

	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
)

type scopeExistFn func(string) (bool, error)

type Role struct {
	id             string
	name           valueobj.SolidString
	description    valueobj.NilString
	assignedScopes *list.List
}

func (r *Role) ChangeDescription(descr string) {
	r.description = valueobj.NewNilString(descr)
}

func (r *Role) AssignScope(scopeId string, existFn scopeExistFn) error {
	scopeToAssign, err := valueobj.NewAssignedScope(r.id, scopeId)
	if err != nil {
		return err
	}

	for elem := r.assignedScopes.Front(); elem != nil; elem = elem.Next() {
		scope, _ := elem.Value.(valueobj.AssignedScope)
		if scope.IsTheSameAs(scopeToAssign) {
			return fmt.Errorf("scope with id %s is already assigned", scopeId)
		}
	}

	if exist, err := existFn(scopeId); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("scope with id %s doesn't exist", scopeId)
	}

	r.assignedScopes.PushBack(scopeToAssign)
	return nil
}

func (r *Role) UnassignScope(scopeId string) error {
	scopeToUnassign, err := valueobj.NewAssignedScope(r.id, scopeId)
	if err != nil {
		return err
	}

	var rmElem *list.Element
	for elem := r.assignedScopes.Front(); elem != nil; elem = elem.Next() {
		scope, _ := elem.Value.(valueobj.AssignedScope)
		if scope.IsTheSameAs(scopeToUnassign) {
			rmElem = elem
			break
		}
	}

	if rmElem == nil {
		return fmt.Errorf("scope with id %s is not assigned", scopeId)
	}

	r.assignedScopes.Remove(rmElem)
	return nil
}

func (r *Role) ToDto() RoleDto {
	return RoleDto{
		Id:          r.id,
		Name:        r.name.String(),
		Description: r.description.Ptr(),
	}
}

func (r *Role) ScopesDto() []AssignedScopeDto {
	dto := make([]AssignedScopeDto, 0)

	for elem := r.assignedScopes.Front(); elem != nil; elem = elem.Next() {
		scope, _ := elem.Value.(valueobj.AssignedScope)
		scopeDto := AssignedScopeDto{RoleId: scope.RoleId(), ScopeId: scope.ScopeId()}
		dto = append(dto, scopeDto)
	}

	return dto
}
