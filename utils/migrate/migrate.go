package main

import "TheBook/sqlite"

func main() {
	err := sqlite.Migrate("thebook.db")
	if err != nil {
		panic(err)
	}
}
