package test

import (
	"os"
	"casino_majiang/conf/config"
	"time"
)

func init() {
	e := config.InitConfig(false)
	if e.IsError() {
		//log.Error("config init failed.", e)
		os.Exit(-1)
	}
	time.Sleep(time.Second * 1)

}