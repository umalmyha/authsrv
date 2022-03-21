package role

import (
	"fmt"

	"container/list"

	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type scopeExistFn func(string) (bool, error)

type Role struct {
	id          string
	name        valueobj.SolidString
	description valueobj.NilString
	scopes      *list.List
}

func (r *Role) ChangeDescription(descr string) {
	r.description = valueobj.NewNilString(descr)
}

func (r *Role) AssignScope(scopeId string, existFn scopeExistFn) error {
	scopeIdent, err := valueobj.NewScopeId(scopeId)
	if err != nil {
		return err
	}

	for elem := r.scopes.Front(); elem != nil; elem = elem.Next() {
		assignedScopeId, _ := elem.Value.(valueobj.ScopeId)
		if assignedScopeId.Equal(scopeIdent) {
			return fmt.Errorf("scope with id %s is already assigned", scopeId)
		}
	}

	if exist, err := existFn(scopeId); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("scope with id %s doesn't exist", scopeId)
	}

	r.scopes.PushBack(scopeIdent)
	return nil
}

func (r *Role) UnassignScope(scopeId string) error {
	scopeIdent, err := valueobj.NewScopeId(scopeId)
	if err != nil {
		return err
	}

	var rmElem *list.Element
	for elem := r.scopes.Front(); elem != nil; elem = elem.Next() {
		assignedScopeId, _ := elem.Value.(valueobj.ScopeId)
		if assignedScopeId.Equal(scopeIdent) {
			rmElem = elem
			break
		}
	}

	if rmElem == nil {
		return fmt.Errorf("scope with id %s is not assigned", scopeId)
	}

	r.scopes.Remove(rmElem)
	return nil
}

func (r *Role) ToDto() RoleDto {
	return RoleDto{
		Id:          r.id,
		Name:        r.name.String(),
		Description: r.description.Ptr(),
	}
}

func (r *Role) ScopesDto() []ScopeAssignmentDto {
	return helpers.FromListWithReducer(r.scopes, func(scopeId valueobj.ScopeId) ScopeAssignmentDto {
		return ScopeAssignmentDto{RoleId: r.id, ScopeId: scopeId.String()}
	})
}
