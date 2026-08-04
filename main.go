package main

import (
	_ "TheBook/config"
	"TheBook/service"
)

func main() {
	server, err := service.InitSystem()
	if err != nil {
		panic(err)
	}

	server.Router.Run(":8080")
}
