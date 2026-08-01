package main

import (
	_ "TheBook/config"
	"TheBook/service"
)

func main() {
	r, _, err := service.InitSystem()
	if err != nil {
		panic(err)
	}

	r.Run(":8080")
}
