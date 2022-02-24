package role

import (
	"fmt"

	"github.com/umalmyha/authsrv/pkg/helpers"
)

type NewRoleDto struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type AssignedScopeDto struct {
	RoleId  string `db:"role_id"`
	ScopeId string `db:"scope_id"`
}

func (dto AssignedScopeDto) Key() string {
	return fmt.Sprintf("%s-%s", dto.RoleId, dto.ScopeId)
}

func (dto AssignedScopeDto) IsPresent() bool {
	return dto.RoleId != "" && dto.ScopeId != ""
}

func (dto AssignedScopeDto) IsTheSameAs(other AssignedScopeDto) bool {
	return dto.RoleId == other.RoleId && dto.ScopeId == other.ScopeId
}

func (dto AssignedScopeDto) Clone() AssignedScopeDto {
	return dto
}

type RoleDto struct {
	Id          string  `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`
}

func (dto RoleDto) Key() string {
	return dto.Id
}

func (dto RoleDto) IsPresent() bool {
	return dto.Id != ""
}

func (dto RoleDto) IsTheSameAs(other RoleDto) bool {
	return dto.Name == other.Name && helpers.EqualValues(dto.Description, other.Description)
}

func (dto RoleDto) Clone() RoleDto {
	return RoleDto{
		Id:          dto.Id,
		Name:        dto.Name,
		Description: helpers.CopyValue(dto.Description),
	}
}
