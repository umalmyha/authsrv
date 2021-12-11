package dto

import (
	"github.com/umalmyha/authsrv/internal/user/store"
	"github.com/umalmyha/authsrv/pkg/database"
)

type User struct {
	Id           string  `json:"id"`
	Username     string  `json:"username"`
	Email        *string `json:"email"`
	PasswordHash string  `json:"-"`
	IsSuperuser  bool    `json:"isSuperuser"`
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	MiddleName   *string `json:"middleName"`
}

func UserDtoFromStore(u store.User) User {
	return User{
		Id:           u.Id,
		Username:     u.Username,
		Email:        database.PtrFromNullString(u.Email),
		PasswordHash: u.PasswordHash,
		FirstName:    database.PtrFromNullString(u.FirstName),
		LastName:     database.PtrFromNullString(u.LastName),
		MiddleName:   database.PtrFromNullString(u.MiddleName),
	}
}
