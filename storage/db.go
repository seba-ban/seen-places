package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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

func (c *DbConnectionConfig) OpenDbConnection() (*pgxpool.Pool, error) {
	ctx := context.Background()
	dbConn, err := pgxpool.New(
		ctx, c.Dsn(),
	)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}
