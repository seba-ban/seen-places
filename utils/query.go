package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type ExpectedQueryType int

const (
	ExpectedQueryTypeInt ExpectedQueryType = iota
	ExpectedQueryTypeFloat
	ExpectedQueryTypeString
	ExpectedQueryTypeRFC3339Time
)

var ErrorQueryNotDefined = errors.New("query not defined")
var ErrorQueryNotValid = errors.New("query not valid")

type QueryValidation struct {
	Name     string
	Type     ExpectedQueryType
	Required bool
}

type QueryConverter struct {
	Ctx *gin.Context
}

func (q *QueryConverter) get(name string) (string, error) {
	val, isSet := q.Ctx.GetQuery(name)
	if !isSet {
		return "", ErrorQueryNotDefined
	}
	return val, nil
}

func (q *QueryConverter) GetRFC3339Time(name string) (*time.Time, error) {
	query, err := q.get(name)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339, query)
	if err != nil {
		return nil, ErrorQueryNotValid
	}
	return &t, nil
}

func (q *QueryConverter) GetInt(name string) (int, error) {
	query, err := q.get(name)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(query)
	if err != nil {
		return 0, ErrorQueryNotValid
	}
	return i, nil
}

func (q *QueryConverter) GetFloat(name string) (float64, error) {
	query, err := q.get(name)
	if err != nil {
		return 0.0, err
	}
	f, err := strconv.ParseFloat(query, 64)
	if err != nil {
		return 0.0, ErrorQueryNotValid
	}
	return f, nil
}

func (q *QueryConverter) GetString(name string) (string, error) {
	return q.get(name)
}

func (q *QueryConverter) GetPgtypeText(name string) (*pgtype.Text, error) {
	query, err := q.get(name)
	if err != nil {
		return nil, err
	}
	return &pgtype.Text{String: query, Valid: true}, nil
}

func (q *QueryConverter) GetPgtypeTime(name string) (*pgtype.Time, error) {
	t, err := q.GetRFC3339Time(name)
	if err != nil {
		return nil, err
	}
	return &pgtype.Time{Microseconds: t.UnixMicro(), Valid: true}, nil
}

func (q *QueryConverter) GetPgtypeTimestamp(name string) (*pgtype.Timestamp, error) {
	t, err := q.GetRFC3339Time(name)
	if err != nil {
		return nil, err
	}
	return &pgtype.Timestamp{Time: *t, Valid: true}, nil
}
