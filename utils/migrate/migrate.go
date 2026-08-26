package main

import "TheBook/sqlite"

func main() {
	err := sqlite.Migrate("your_database.db")
	if err != nil {
		panic(err)
	}
}
