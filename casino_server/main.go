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
	"casino_server/bonus"
	"casino_server/system/casinoWeb"
)

func init() {
	fmt.Println("1,runtime.GOMAXPROCS()...")
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("2,InitConfig()...")
	e := config.InitConfig(false)
	if e.IsError() {
		log.Error("config init failed.", e)
		os.Exit(-1)
	}
	fmt.Println("3,initSeed...")
	s := time.Now().UTC().UnixNano()
	rand.Seed(s)
	fmt.Println("4,config init ok...")
}

func main() {
	log.T("main start...")
	lconf.LogLevel = conf.Server.LogLevel        //通过 conf/server.json 去初始化conf.Server
	lconf.LogPath = conf.Server.LogPath        //conf.Server.LogPath

	//后台管理
	go func() {
		log.T("web start...")
		casinoWeb.InitCms()
	}()

	//初始化三个模块,主函数入口
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
		bonus.Module,
	)

}
