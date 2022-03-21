package user

import (
	"fmt"

	"github.com/umalmyha/authsrv/pkg/helpers"
)

type UserDto struct {
	Id          string  `db:"id"`
	Username    string  `db:"username"`
	Email       *string `db:"email"`
	Password    string  `db:"password_hash"`
	IsSuperuser bool    `db:"is_superuser"`
	FirstName   *string `db:"first_name"`
	LastName    *string `db:"last_name"`
	MiddleName  *string `db:"middle_name"`
}

func (dto UserDto) Key() string {
	return dto.Id
}

func (dto UserDto) IsPresent() bool {
	return dto.Id != ""
}

func (dto UserDto) Equal(other UserDto) bool {
	return dto.Username == other.Username &&
		dto.Password == other.Password &&
		dto.IsSuperuser == other.IsSuperuser &&
		helpers.EqualValues(dto.Email, other.Email) &&
		helpers.EqualValues(dto.FirstName, other.FirstName) &&
		helpers.EqualValues(dto.LastName, other.LastName) &&
		helpers.EqualValues(dto.MiddleName, other.MiddleName)
}

func (dto UserDto) Clone() UserDto {
	return UserDto{
		Id:          dto.Id,
		Username:    dto.Username,
		IsSuperuser: dto.IsSuperuser,
		Password:    dto.Password,
		Email:       helpers.CopyValue(dto.Email),
		FirstName:   helpers.CopyValue(dto.FirstName),
		LastName:    helpers.CopyValue(dto.LastName),
		MiddleName:  helpers.CopyValue(dto.MiddleName),
	}
}

type RoleAssignmentDto struct {
	UserId string `db:"user_id"`
	RoleId string `db:"role_id"`
}

func (dto RoleAssignmentDto) Key() string {
	return fmt.Sprintf("%s-%s", dto.UserId, dto.RoleId)
}

func (dto RoleAssignmentDto) IsPresent() bool {
	return dto.UserId != "" && dto.RoleId != ""
}

func (dto RoleAssignmentDto) Equal(other RoleAssignmentDto) bool {
	return dto.UserId == other.UserId && dto.RoleId == other.RoleId
}

func (dto RoleAssignmentDto) Clone() RoleAssignmentDto {
	return dto
}

type NewUserDto struct {
	Username        string  `json:"username"`
	Email           *string `json:"email"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirmPassword"`
	FirstName       *string `json:"firstName"`
	LastName        *string `json:"lastName"`
	MiddleName      *string `json:"middleName"`
	IsSuperuser     bool    `json:"-"`
}

type UserAuthDto struct {
	UserId    string `db:"user_id"`
	RoleId    string `db:"role_id"`
	RoleName  string `db:"role_name"`
	ScopeId   string `db:"scope_id"`
	ScopeName string `db:"scope_name"`
}

type SigninDto struct {
	Username    string `json:"user"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type LogoutDto struct {
	Username       string `json:"user"`
	Fingerprint    string `json:"fingerprint"`
	RefreshTokenId string `json:"-"`
}

type RefreshDto struct {
	Username       string `json:"user"`
	Fingerprint    string `json:"fingerprint"`
	RefreshTokenId string `json:"-"`
}
