package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DbConnectionConfig struct {
	Host     string
	Port     int
	DbName   string
	User     string
	Password string
	Sslmode  string
}

func (c *DbConnectionConfig) Dsn() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DbName, c.Sslmode,
	)
}

func (c *DbConnectionConfig) OpenDbConnection() (*context.Context, *pgx.Conn, error) {
	ctx := context.Background()
	dbConn, err := pgx.Connect(
		ctx, c.Dsn(),
	)
	if err != nil {
		return nil, nil, err
	}

	return &ctx, dbConn, nil
}
