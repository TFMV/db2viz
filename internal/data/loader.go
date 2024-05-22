package data

import (
	"context"
	"log"

	"github.com/TFMV/db2viz/internal/db"
)

type Loader struct {
	dbConn *db.PostgresConnection
	table  string
}

func NewLoader(dbConn *db.PostgresConnection, table string) *Loader {
	return &Loader{
		dbConn: dbConn,
		table:  table,
	}
}

func (l *Loader) LoadData(ctx context.Context) ([]map[string]interface{}, error) {
	log.Printf("Loading data from table: %s", l.table)
	return db.LoadData(l.dbConn, l.table)
}
