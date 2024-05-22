package data

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Loader struct {
	conn *pgx.Conn
}

func NewLoader(conn *pgx.Conn) *Loader {
	return &Loader{conn: conn}
}

func (l *Loader) LoadData(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := l.conn.Query(ctx, "SELECT * FROM your_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		record := make(map[string]interface{})
		for i, col := range rows.FieldDescriptions() {
			record[string(col.Name)] = values[i]
		}
		records = append(records, record)
	}

	return records, nil
}
