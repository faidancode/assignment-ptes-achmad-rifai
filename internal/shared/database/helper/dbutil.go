package helper

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

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

// BoolValue mengonversi *bool ke bool dengan fallback nilai default jika nil.
func BoolValue(b *bool, defaultValue bool) bool {
	if b == nil {
		return defaultValue
	}
	return *b
}

// NewNullBool membantu konversi *bool -> sql.NullBool untuk keperluan database
func NewNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  *b,
		Valid: true,
	}
}

// ToDecimal membantu konversi float64 ke decimal.Decimal untuk sqlc
func ToDecimal(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

// FloatFromDecimal membantu konversi balik dari database (decimal) ke response (float64)
func FloatFromDecimal(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

// NewNullDecimal menangani pemetaan harga (decimal) opsional
func NewNullDecimal(f *float64) decimal.NullDecimal {
	if f == nil {
		return decimal.NullDecimal{Valid: false}
	}
	// Konversi float64 ke decimal.Decimal
	d := decimal.NewFromFloat(*f)
	return decimal.NullDecimal{Decimal: d, Valid: true}
}

// NewNullInt32 menangani pemetaan stock (int32) opsional
func NewNullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}
