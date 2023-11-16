package utils

import (
	"time"

	"github.com/hako/durafmt"
	"github.com/jackc/pgx/v5/pgtype"
)

func getTime(t any) time.Time {
	switch v := t.(type) {
	case time.Time:
		return v
	case pgtype.Timestamp:
		return v.Time
	default:
		panic("unknown time type")
	}
}

func FormatTime(t any) string {
	return getTime(t).Format(time.RFC3339)
}

func FormatTimeDuration(start any, end any) string {
	timeduration := getTime(end).Sub(getTime(start))
	return durafmt.Parse(timeduration).String()
}
