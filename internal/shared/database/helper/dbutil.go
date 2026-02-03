package helper

import (
	"database/sql"
	"strconv"

	"github.com/shopspring/decimal"
)

//
// =======================
// STRING
// =======================
//

func StringValue(s string) string {
	return s
}

func StringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func StringToNull(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

//
// =======================
// BOOL
// =======================
//

func BoolValue(b bool) bool {
	return b
}

func BoolPtrValue(b *bool, defaultValue bool) bool {
	if b == nil {
		return defaultValue
	}
	return *b
}

func BoolToNull(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

//
// =======================
// INT32
// =======================
//

func Int32Value(i int32) int32 {
	return i
}

func Int32PtrValue(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func Int32ToNull(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

//
// =======================
// FLOAT64
// =======================
//

func Float64Value(f float64) float64 {
	return f
}

func Float64PtrValue(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

//
// =======================
// DECIMAL
// =======================
//

func DecimalValue(d decimal.Decimal) decimal.Decimal {
	return d
}

func DecimalPtrValue(d *decimal.Decimal) decimal.Decimal {
	if d == nil {
		return decimal.Zero
	}
	return *d
}

func Float64ToDecimal(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}

func Float64PtrToDecimal(f *float64) decimal.Decimal {
	if f == nil {
		return decimal.Zero
	}
	return decimal.NewFromFloat(*f)
}

// Presisi tinggi (uang / harga)
func Float64ToDecimalExact(f float64) decimal.Decimal {
	return decimal.RequireFromString(
		strconv.FormatFloat(f, 'f', -1, 64),
	)
}

func Float64PtrToDecimalExact(f *float64) decimal.Decimal {
	if f == nil {
		return decimal.Zero
	}
	return Float64ToDecimalExact(*f)
}

func DecimalToFloat64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

//
// =======================
// DECIMAL → NULL
// =======================
//

func Float64ToNullDecimal(f *float64) decimal.NullDecimal {
	if f == nil {
		return decimal.NullDecimal{}
	}
	return decimal.NullDecimal{
		Decimal: decimal.NewFromFloat(*f),
		Valid:   true,
	}
}
