package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"runtime"
	"casino_login/conf"
	"casino_login/game"
	"casino_login/gate"
	"casino_login/login"
)

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
