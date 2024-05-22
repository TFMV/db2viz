package data

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Loader struct {
	db *pgxpool.Pool
}

func NewLoader(db *pgxpool.Pool) *Loader {
	return &Loader{db: db}
}

func (l *Loader) LoadData(ctx context.Context, table string) ([]map[string]interface{}, error) {
	rows, err := l.db.Query(ctx, "SELECT * FROM "+table)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[string(col.Name)] = values[i]
		}
		data = append(data, rowMap)
		// log.Printf("Loaded row: %+v\n", rowMap)
	}
	return data, nil
}
