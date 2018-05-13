package main

import (
	"flag"
	"github.com/therako/sql-exporter/exporter"
	"os"
)

func main() {
	var config exporter.Config
	flag.StringVar(&config.SqlBind, "conn", "", "Text to parse. (Required)")
	flag.StringVar(&config.SqlDriver, "driver", "postgres", "Metric {postgres|mysql};. (Default: postgres)")
	flag.StringVar(&config.Query, "query", "", "Sql Query. (Required)")
	flag.StringVar(&config.OutputFile, "output", "./tmp", "Full output path with prefix if any (Default: ./tmp)")
	flag.IntVar(&config.Concurrency, "concurrency", 1, "no. of concurrency workers (Default: 1)")
	flag.Parse()

	if config.SqlBind == "" || config.Query == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if config.Concurrency <= 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	exporter.Export(config)
}
