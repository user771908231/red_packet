package test

import (
	"runtime"
	"fmt"
	"os"
	"time"
	"casino_server/common/config"
	"casino_server/common/log"
	"math/rand"
)

func init(){
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
