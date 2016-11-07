package test

import (
	"os"
	"time"
	"casino_doudizhu/conf/config"
)

func init() {
	e := config.InitConfig(false)
	if e.IsError() {
		//log.Error("config init failed.", e)
		os.Exit(-1)
	}
	time.Sleep(time.Second * 1)
}

