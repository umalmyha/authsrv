package store

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/umalmyha/authsrv/pkg/database"
)

type Store struct {
	*database.Store
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		Store: database.NewStore(db),
	}
}

func (s *Store) Create(ctx context.Context, u User) database.SqlContextExecFunc {
	return func(ec sqlx.ExtContext) error {
		q := "INSERT INTO USERS(ID, USERNAME, EMAIL, PASSWORD_HASH, FIRST_NAME, LAST_NAME, MIDDLE_NAME, IS_SUPERUSER) VALUES($1, $2, $3, $4, $5, $6, $7, $8)"
		_, err := ec.ExecContext(ctx, q, u.Id, u.Username, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.MiddleName, u.IsSuperuser)
		if err != nil {
			return err
		}
		return nil
	}
}

func (s *Store) All(ctx context.Context, dest *[]User) database.SqlContextExecFunc {
	return func(ec sqlx.ExtContext) error {
		q := "SELECT * FROM USERS"
		return sqlx.SelectContext(ctx, ec, dest, q)
	}
}

func (s *Store) ById(ctx context.Context, dest *User, id string) database.SqlContextExecFunc {
	return func(ec sqlx.ExtContext) error {
		q := "SELECT * FROM USERS WHERE ID = $1 LIMIT 1"
		return sqlx.GetContext(ctx, ec, dest, q, id)
	}
}

func (s *Store) Update(ctx context.Context, u User) database.SqlContextExecFunc {
	return func(ec sqlx.ExtContext) error {
		q := "UPDATE USER SET EMAIL = $1, FIRST_NAME = $2, LAST_NAME = $3, MIDDLE_NAME = $4"
		if _, err := ec.ExecContext(ctx, q, u.Email, u.FirstName, u.LastName, u.MiddleName); err != nil {
			return err
		}
		return nil
	}
}

func (s *Store) Delete(ctx context.Context, id string) database.SqlContextExecFunc {
	return func(ec sqlx.ExtContext) error {
		q := "DELETE FROM USERS WHERE ID = $1"
		if _, err := ec.ExecContext(ctx, q, id); err != nil {
			return err
		}
		return nil
	}
}

func (s *Store) ByUsername(ctx context.Context, dest *User, username string) database.SqlContextExecFunc {
	return func(ec sqlx.ExtContext) error {
		q := "SELECT * FROM USERS WHERE USERNAME = $1 LIMIT 1"
		if err := sqlx.GetContext(ctx, ec, dest, q, username); err != nil {
			return err
		}
		return nil
	}
}
