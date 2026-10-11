package main

import (
	_ "TheBook/config"
	"TheBook/logger"
	"TheBook/service"
	"TheBook/utils"
	"fmt"
)

// main 初始化日志和服务，并启动本地 HTTP 服务器。
func main() {
	rootPath, err := utils.FindProjectRoot()
	if err != nil {
		panic(err)
	}

	log, closeLog, err := logger.InitLogger(fmt.Sprintf("%s/logs/app.log", rootPath), true)
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
