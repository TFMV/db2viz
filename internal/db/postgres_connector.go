package db

import (
	"context"
	"fmt"

	"github.com/TFMV/db2viz/config"
	"github.com/jackc/pgx/v5"
)

func ConnectPostgres(cfg config.PostgresConfig) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	conn, err := pgx.Connect(context.Background(), connStr)
	return conn, err
}
