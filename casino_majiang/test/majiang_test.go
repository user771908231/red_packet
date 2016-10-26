package test

import (
	"testing"
	"casino_majiang/conf/config"
	"casino_majiang/gamedata/dao"
	"fmt"
)

func Test(t *testing.T) {
	config.InitConfig(false)
	data := dao.GetByUserId(501)
	fmt.Println("找到的数据", data)
}
