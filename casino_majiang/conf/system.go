package conf

import "github.com/name5566/leaf/log"

func InitSystem() {
	initMongo()
	initRedis()
}

func initMongo() {
	log.Debug("初始化mongo")

}

func initRedis() {
	log.Debug("初始化redis")

}