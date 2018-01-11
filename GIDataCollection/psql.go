package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Record map[string]interface{}
type DataSet []Record

type Postgres struct {
	db *sql.DB
}

func (pg *Postgres) Open() (err error) {
	if pg.db != nil {
		pg.Close()
	}

	connStr := fmt.Sprintf("port=%d dbname=%s user=%s password=%s sslmode=disable", DBPort, DBName, DBUser, DBPassword)
	pg.db, err = sql.Open("postgres", connStr)
	return
}

func (pg *Postgres) Close() {
	pg.db.Close()
}

func (pg *Postgres) GetDB() (*sql.DB) {
	return pg.db
}

func (pg *Postgres) Query(sql string, args ...interface{}) (DataSet, error) {
	rows, err := pg.db.Query(sql, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	colNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	colCount := len(colNames)
	container := make([]interface{}, colCount)
	vars := make([]interface{}, colCount)
	for index := range container  {
		container[index] = &vars[index]
	}

	dataSet := DataSet{}
	for rows.Next()  {
		err = rows.Scan(container...)
		if err != nil {
			return nil, err
		}

		var rec = Record{}
		for index, colName := range colNames {
			rec[colName] = vars[index]
		}

		dataSet = append(dataSet, rec)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return dataSet, nil
}

func (pg *Postgres) Exec(sql string, args ...interface{}) (int64, error) {
	result, err := pg.db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
