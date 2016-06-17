package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"casino_server/game"
	"casino_server/conf"
	"casino_server/gate"
	"casino_server/login"
	"runtime"
	"fmt"
	"os"
	"time"
	"casino_server/common/config"
	"casino_server/common/log"
	"math/rand"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("bbsvr >>>> config init...")

	e := config.InitConfig(false)
	if e.IsError() {
		log.Error("config init failed.", e)
		os.Exit(-1)
	}
	log.Normal("config init ok...")

	//随机种子
	s := time.Now().UTC().UnixNano()
	log.Normal("Server start... seed: %v", s)
	rand.Seed(s)
}



func main() {
	log.T("main start...")
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
