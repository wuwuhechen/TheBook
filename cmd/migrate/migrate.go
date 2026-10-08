package main

import (
	"flag"
	"log"

	"TheBook/sqlite"
)

func main() {
	databasePath := flag.String(
		"database",
		"database/test.db",
		"SQLite 数据库路径",
	)

	flag.Parse()

	if err := sqlite.Migrate(*databasePath); err != nil {
		log.Fatal(err)
	}
}
