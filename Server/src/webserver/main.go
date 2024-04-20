package main

import(
	"common/config"
	"common/connection"
	"common/model"
	"common/redis"
	"common/manager"
	"common/log"
	"webserver/router"
)

func main() {
	config.InitGlobalConfig()
	connection.InitMySQL()
	model.InitModel()
	manager.InitManager()
	redis.InitRedis()
	log.Init()
	router.InitRouter()

	router.Start()
}