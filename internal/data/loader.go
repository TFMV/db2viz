package data

import (
	"context"
	"log"

	"github.com/TFMV/db2viz/internal/db"
)

type Loader struct {
	dbConn *db.PostgresConnection
	schema string
	table  string
}

func NewLoader(dbConn *db.PostgresConnection, schema string, table string) *Loader {
	return &Loader{
		dbConn: dbConn,
		schema: schema,
		table:  table,
	}
}

func (l *Loader) LoadData(ctx context.Context) ([]map[string]interface{}, error) {
	log.Printf("Loading data from table: %s.%s", l.schema, l.table)
	return db.LoadData(l.dbConn, l.schema, l.table)
}
