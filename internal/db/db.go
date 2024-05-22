package db

import (
	"context"
	"fmt"
	"log"

	"github.com/TFMV/db2viz/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConnection struct {
	Pool *pgxpool.Pool
}

func NewPostgresConnection(cfg config.PostgresConfig) (*PostgresConnection, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to PostgreSQL successfully")
	return &PostgresConnection{Pool: pool}, nil
}

func (pc *PostgresConnection) Close(ctx context.Context) {
	pc.Pool.Close()
}

func LoadData(dbConn *PostgresConnection, table string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := dbConn.Pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns := rows.FieldDescriptions()
	var data []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[string(col.Name)] = values[i]
		}
		data = append(data, rowMap)
	}
	return data, nil
}
