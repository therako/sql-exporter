package main

import (
	"flag"
	"os"
)

func main() {
	var config Config
	flag.StringVar(&config.SqlBind, "conn", "", "Text to parse. (Required)")
	flag.StringVar(&config.SqlDriver, "driver", "postgres", "Metric {postgres|mysql};. (Default: postgres)")
	flag.StringVar(&config.Query, "query", "", "Sql Query. (Required)")
	flag.StringVar(&config.OutputFile, "output", "", "Export file full path. (Required)")
	flag.IntVar(&config.Concurrency, "concurrency", 1, "no. of concurrency workers (Default: 1)")
	flag.Parse()

	if config.SqlBind == "" || config.Query == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	Export(config)
}
