package services

import "database/sql"

func processNullOrString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}