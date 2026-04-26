package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Pgtype Int4 -> Int64
func ConvertInt4ToInt64(value pgtype.Int4) int64 {
	if value.Valid {
		return int64(value.Int32)
	}
	return 0
}

// Pgtype Int8 -> Int64
func ConvertInt8ToInt64(value pgtype.Int8) int64 {
	if value.Valid {
		return value.Int64
	}
	return 0
}

// PgtypeBool -> Bool
func ConvertPgtypeBoolToBool(value pgtype.Bool) bool {
	if value.Valid {
		return value.Bool
	}
	return false
}

// Bool -> PgtypeBool
func ConvertBoolToPgtypeBool(value bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  value,
		Valid: true,
	}
}

// String -> PgtypeText
func ConvertStringToPgtypeText(input *string) pgtype.Text {
	var value pgtype.Text
	if input == nil || *input == "" {
		_ = value.Scan(nil)
	} else {
		_ = value.Scan(*input)
	}
	return value
}

// PgtypeText -> String
func ConvertPgtypeTextToString(value pgtype.Text) string {
	if value.Valid {
		return value.String
	}
	return ""
}

// String -> PgtypeDate
func ConvertStringToPgtypeDate(input *string) pgtype.Date {
	if input == nil || *input == "" {
		return pgtype.Date{
			Valid: false,
		}
	}
	t, err := time.Parse("2006-01-02", *input)
	if err != nil {
		return pgtype.Date{
			Valid: false,
		}
	}
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

// Pgtype.Date -> String
func ConvertPgtypeDateToString(value pgtype.Date) string {
	if value.Valid {
		return value.Time.Format("2006-01-02")
	}
	return ""
}

// Time -> PgtypeTimestamptz
func ConvertTimeToPgtypeTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

// PgtypeTimestamptz -> Timestamp
func ConvertPgtypeTimestamptzToTimestamp(value pgtype.Timestamptz) *timestamppb.Timestamp {
	if !value.Valid || value.Time.IsZero() {
		return nil
	}
	return timestamppb.New(value.Time)
}
