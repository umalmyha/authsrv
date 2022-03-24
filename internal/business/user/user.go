package user

import (
	"container/list"
	"errors"
	"fmt"
	"time"

	"github.com/umalmyha/authsrv/internal/business/refresh"
	"github.com/umalmyha/authsrv/internal/business/role"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/helpers"
)

type roleExistFn func(string) (bool, error)

type User struct {
	id          string
	username    valueobj.SolidString
	email       valueobj.NilEmail
	password    valueobj.Password
	isSuperuser bool
	firstName   valueobj.NilString
	lastName    valueobj.NilString
	middleName  valueobj.NilString
	roles       *list.List
	tokens      *list.List
	auth        valueobj.UserAuth
}

type RoleFinderByNameFn func(string) (role.RoleDto, error)

func (u *User) AssignRole(name string, finderFn RoleFinderByNameFn) error {
	r, err := finderFn(name)
	if err != nil {
		return err
	}

	if !r.IsPresent() {
		return fmt.Errorf("role %s doesn't exist", name)
	}

	roleIdent, err := valueobj.NewRoleId(r.Id)
	if err != nil {
		return err
	}

	for elem := u.roles.Front(); elem != nil; elem = elem.Next() {
		assignedRoleId, _ := elem.Value.(valueobj.RoleId)
		if assignedRoleId.Equal(roleIdent) {
			return fmt.Errorf("role %s is already assigned", name)
		}
	}

	u.roles.PushBack(roleIdent)
	return nil
}

func (u *User) UnassignRole(name string, finderFn RoleFinderByNameFn) error {
	r, err := finderFn(name)
	if err != nil {
		return err
	}

	if !r.IsPresent() {
		return fmt.Errorf("role %s doesn't exist", name)
	}

	roleIdent, err := valueobj.NewRoleId(r.Id)
	if err != nil {
		return err
	}

	var rmElem *list.Element
	for elem := u.roles.Front(); elem != nil; elem = elem.Next() {
		assignedRoleId, _ := elem.Value.(valueobj.RoleId)
		if assignedRoleId.Equal(roleIdent) {
			rmElem = elem
			break
		}
	}

	if rmElem == nil {
		return fmt.Errorf("role %s is not assigned to user %s", name, u.username)
	}

	u.roles.Remove(rmElem)
	return nil
}

func (u *User) GenerateJwt(issuedAt time.Time, cfg valueobj.JwtConfig) (valueobj.Jwt, error) {
	return valueobj.NewJwt(u.username.String(), issuedAt, u.auth.Roles(), u.auth.Scopes(), cfg)
}

func (u *User) GenerateRefreshToken(fgrprint string, issuedAt time.Time, cfg valueobj.RefreshTokenConfig) (*refresh.RefreshToken, error) {
	for elem := u.tokens.Front(); elem != nil; elem = elem.Next() {
		token, _ := elem.Value.(*refresh.RefreshToken)
		if token.Fingerprint() == fgrprint {
			return nil, errors.New(fmt.Sprintf("refresh token for device %s is generated already", fgrprint))
		}
	}

	token, err := refresh.NewRefreshToken(fgrprint, issuedAt, cfg)
	if err != nil {
		return nil, err
	}

	if u.tokens.Len() == cfg.MaxTokensCount() {
		u.removeTokenClosestToExpiration()
	}
	u.tokens.PushBack(token)

	return token, nil
}

func (u *User) DiscardRefreshToken(logout LogoutDto) error {
	if logout.RefreshTokenId == "" {
		return errors.New("refresh token id can't be initial")
	}

	if logout.Fingerprint == "" {
		return errors.New("fingerprint must be provided")
	}

	rmElem := u.findRefreshTokenElemById(logout.RefreshTokenId)
	if rmElem == nil {
		return errors.New(fmt.Sprintf("provided refresh token doesn't exist or doesn't belong to user %s", u.username))
	}

	token, _ := rmElem.Value.(*refresh.RefreshToken)
	if token.Fingerprint() != logout.Fingerprint {
		return errors.New("provided fingerprint doesn't belong to provided refresh token")
	}

	u.tokens.Remove(rmElem)
	return nil
}

func (u *User) RefreshSession(rfr RefreshDto, now time.Time) error {
	// TODO: think of reusage of error handling
	if rfr.RefreshTokenId == "" {
		return errors.New("refresh token id can't be initial")
	}

	if rfr.Fingerprint == "" {
		return errors.New("fingerprint must be provided")
	}

	tokenElem := u.findRefreshTokenElemById(rfr.RefreshTokenId)
	if tokenElem == nil {
		return errors.New(fmt.Sprintf("provided refresh token doesn't exist or doesn't belong to user %s", u.username))
	}
	u.tokens.Remove(tokenElem)

	token, _ := tokenElem.Value.(*refresh.RefreshToken)
	if token.Fingerprint() != rfr.RefreshTokenId {
		return errors.New("provided refresh token doesn't exist or doesn't belong to user %s")
	}

	if token.ExpiresAt().Before(now) {
		return errors.New("provided refresh token doesn't exist or doesn't belong to user %s")
	}

	return nil
}

func (u *User) VerifyPassword(password string) (bool, error) {
	if password == "" {
		return false, errors.New("password for verification can't be initial")
	}

	hash, err := valueobj.GenerateHash([]byte(password))
	if err != nil {
		return false, err
	}

	return hash == u.password.Hash(), nil
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

func (u *User) RolesDto() []RoleAssignmentDto {
	return helpers.FromListWithReducer(u.roles, func(roleId valueobj.RoleId) RoleAssignmentDto {
		return RoleAssignmentDto{UserId: u.id, RoleId: roleId.String()}
	})
}

func (u *User) TokensDto() []refresh.RefreshTokenDto {
	return helpers.FromListWithReducer(u.tokens, func(token *refresh.RefreshToken) refresh.RefreshTokenDto {
		return refresh.RefreshTokenDto{
			Id:          token.Id(),
			Fingerprint: token.Fingerprint(),
			UserId:      u.id,
			IssuedAt:    token.IssuedAt(),
			ExpiresAt:   token.ExpiresAt(),
		}
	})
}

func (u *User) removeTokenClosestToExpiration() {
	var rmElem *list.Element

	tokenExpiresEarlier := func(rm *list.Element, curr *list.Element) bool {
		rmToken, _ := rm.Value.(*refresh.RefreshToken)
		currToken, _ := curr.Value.(*refresh.RefreshToken)
		return currToken.ExpiresAt().Before(rmToken.ExpiresAt())
	}

	for elem := u.tokens.Front(); elem != nil; elem = elem.Next() {
		if rmElem == nil {
			rmElem = elem
			continue
		}

		if tokenExpiresEarlier(rmElem, elem) {
			rmElem = elem
		}
	}

	u.tokens.Remove(rmElem)
}

func (u *User) findRefreshTokenElemById(tokenId string) *list.Element {
	for elem := u.tokens.Front(); elem != nil; elem = elem.Next() {
		token, _ := elem.Value.(*refresh.RefreshToken)
		if token.Id() == tokenId {
			return elem
		}
	}
	return nil
}
