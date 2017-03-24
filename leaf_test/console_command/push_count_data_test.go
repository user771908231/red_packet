package console_command

import (
	"casino_common/common/service/pushService"
	"testing"
	"casino_common/common/service/countService"
	"math/big"
)

func TestPushCountData(t *testing.T) {
	pushService.PoolInit("192.168.199.200:2801")
	pushService.PushUserData(11750)
	row := countService.T_game_log{
		UserId: 13801,
		IsWine:true,
		Bill: 9595,
	}
	t.Log(row.Insert())
	t.Log(row.Insert())
}

func TestFloatAdd(t *testing.T) {
	var f1 float64 = 0.02
	var f2 float64 = 0.099
	f3 := f1 + f2

	//保留两位有效数字
	f3 = float64(int64(f3*100))/100

	big1 := big.NewFloat(f1)
	big2 := big.NewFloat(f2)
	big3,_ := big.NewFloat(0).Add(big1, big2).Float64()

	t.Log(f1, f2)
	t.Log(f3, big3)
}
