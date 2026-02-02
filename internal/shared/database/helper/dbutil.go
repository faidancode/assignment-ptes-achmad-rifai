package helper

import "database/sql"

func strPtr(s string) *string {
	return &s
}

// NewNullString membantu konversi *string -> sql.NullString
func NewNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
