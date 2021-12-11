package store

import (
	"database/sql"
)

type User struct {
	Id           string         `db:"id"`
	Username     string         `db:"username"`
	Email        sql.NullString `db:"email"`
	PasswordHash string         `db:"password_hash"`
	IsSuperuser  bool           `db:"is_superuser"`
	FirstName    sql.NullString `db:"first_name"`
	LastName     sql.NullString `db:"last_name"`
	MiddleName   sql.NullString `db:"middle_name"`
}
