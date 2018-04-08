package main

type Config struct {
	SqlBind     string
	SqlDriver   string
	Query       string
	OutputFile  string
	Concurrency int
}
