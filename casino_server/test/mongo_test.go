package test

import (
	"testing"
	"fmt"
	"casino_server/mode/dao/TThDeskRecordDao"
)

func TestMongoUtils(t *testing.T) {
	fmt.Println("开始测试")
	var userId uint32 = 10180
	d2 := TThDeskRecordDao.Find(userId,20)
	fmt.Println("d2:", d2)

}

//找有某个人的战绩
