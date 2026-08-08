package main

import (
	_ "TheBook/config"
	"TheBook/logger"
	"TheBook/service"
)

func main() {
	log, closeLog, err := logger.InitLogger("logs/app.log", true)
	if err != nil {
		panic(err)
	}

	defer closeLog()

	server, err := service.InitSystem(log)
	if err != nil {
		panic(err)
	}

	server.Router.Run(":8080")
}
