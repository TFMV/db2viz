package db

import (
	"database/sql"
	"log"
)

type DB2Connection struct {
	*sql.DB
}

func NewDB2Connection(cfg struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}) (*DB2Connection, error) {
	dsn := "HOSTNAME=" + cfg.Host + ";PORT=" + cfg.Port + ";DATABASE=" + cfg.DBName + ";UID=" + cfg.Username + ";PWD=" + cfg.Password + ";"
	db, err := sql.Open("go_ibm_db", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Connected to DB2 successfully")
	return &DB2Connection{db}, nil
}

func LoadData(dbConn *DB2Connection) ([]map[string]interface{}, error) {
	rows, err := dbConn.Query("SELECT * FROM your_table")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

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
			rowMap[col] = values[i]
		}
		data = append(data, rowMap)
	}
	return data, nil
}
