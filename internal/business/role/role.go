package role

import (
	"container/list"

	"github.com/pkg/errors"
	"github.com/umalmyha/authsrv/internal/business/scope"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type ScopeFinderByNameFn func(string) (scope.ScopeDto, error)

type Role struct {
	id          string
	name        valueobj.SolidString
	description valueobj.NilString
	scopes      *list.List
}

func (r *Role) ChangeDescription(descr string) {
	r.description = valueobj.NewNilString(descr)
}

func (r *Role) AssignScope(name string, finderFn ScopeFinderByNameFn) error {
	sc, err := finderFn(name)
	if err != nil {
		return errors.Wrap(err, "failed to find scope")
	}

	if !sc.IsPresent() {
		return errors.Errorf("scope %s doesn't exist", name)
	}

	scopeIdent, err := valueobj.NewScopeId(sc.Id)
	if err != nil {
		return errors.Wrap(err, "failed to build scope identifier")
	}

	for elem := r.scopes.Front(); elem != nil; elem = elem.Next() {
		assignedScopeId, _ := elem.Value.(valueobj.ScopeId)
		if assignedScopeId.Equal(scopeIdent) {
			return errors.Errorf("scope %s is already assigned", name)
		}
	}

	r.scopes.PushBack(scopeIdent)
	return nil
}

func (r *Role) UnassignScope(name string, finderFn ScopeFinderByNameFn) error {
	sc, err := finderFn(name)
	if err != nil {
		return err
	}

	if !sc.IsPresent() {
		return errors.Errorf("scope %s doesn't exist", name)
	}

	scopeIdent, err := valueobj.NewScopeId(sc.Id)
	if err != nil {
		return errors.Wrap(err, "failed to build scope identifier")
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
		return errors.Errorf("scope %s is not assigned to role %s", name, r.name)
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
