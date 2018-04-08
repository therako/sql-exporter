package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"sync"
)

type RowValues struct {
	RowData []interface{}
}

func Export(config Config) {
	db, err := sql.Open(config.SqlDriver, config.SqlBind)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(config.Query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var cols = make([]string, len(columns))
	copy(cols, columns)
	rowValues := RowValues{}
	rowValues.RowData = make([]interface{}, len(columns))
	for i, _ := range columns {
		rowValues.RowData[i] = &columns[i]
	}

	chanVals := make(chan RowValues)
	var w sync.WaitGroup
	w.Add(config.Concurrency)
	for i := 1; i <= config.Concurrency; i++ {
		file, err := os.Create(fmt.Sprintf("%s_%v.txt", config.OutputFile, i))
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer file.Close()
		go func(ci chan RowValues, filePointer *os.File) {
			for rows := range ci {
				var m = map[string]interface{}{}
				for i, col := range cols {
					m[col] = rows.RowData[i]
				}
				obj, _ := json.Marshal(m)
				filePointer.Write(obj)
				filePointer.Write([]byte("\n"))
			}
		}(chanVals, file)
	}

	for rows.Next() {
		rows.Scan(rowValues.RowData...)
		chanVals <- rowValues
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= noWorker; i++ {
		w.Done()
	}
	close(chanVals)
	w.Wait()
}
