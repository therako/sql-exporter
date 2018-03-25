package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	var config Config
	flag.StringVar(&config.SqlBind, "conn", "", "Text to parse. (Required)")
	flag.StringVar(&config.SqlDriver, "driver", "postgres", "Metric {postgres|mysql};. (Default: postgres)")
	flag.StringVar(&config.Query, "query", "", "Sql Query. (Required)")
	flag.StringVar(&config.OutputFile, "output", "", "Export file full path. (Required)")
	flag.Parse()

	if config.SqlBind == "" || config.Query == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Conn: %s, Driver: %s, Query: %s\n", config.SqlBind, config.SqlDriver, config.Query)

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
	var vals = make([]interface{}, len(columns))
	for i, _ := range columns {
		vals[i] = &columns[i]
	}

	file, err := os.Create(config.OutputFile)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	for rows.Next() {
		err := rows.Scan(vals...)
		if err != nil {
			fmt.Println(err)
		}
		var m = map[string]interface{}{}
		for i, col := range cols {
			m[col] = vals[i]
		}
		obj, _ := json.Marshal(m)
		file.Write(obj)
		file.Write([]byte("\n"))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
