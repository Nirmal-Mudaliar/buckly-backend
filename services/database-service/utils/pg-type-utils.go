package utils

import "github.com/jackc/pgx/v5/pgtype"

func ConvertInt4ToInt64(value pgtype.Int4) int64 {
	if value.Valid {
		return int64(value.Int32)
	}
	return 0
}

func ConvertInt8ToInt64(value pgtype.Int8) int64 {
	if value.Valid {
		return value.Int64
	}
	return 0
}

func ConvertPgtypeTextToString(value pgtype.Text) string {
	if value.Valid {
		return value.String
	}
	return ""
}

func ConvertPgtypeBoolToBool(value pgtype.Bool) bool {
	if value.Valid {
		return value.Bool
	}
	return false
}
