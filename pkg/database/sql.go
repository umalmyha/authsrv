package database

import "database/sql"

func NewNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func PtrFromNullString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}
