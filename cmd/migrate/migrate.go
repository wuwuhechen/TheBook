package main

import (
	"flag"
	"log"

	"TheBook/sqlite"
)

// main 解析数据库文件名并执行数据迁移命令。
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
