package console_command

import (
	"casino_common/common/service/pushService"
	"testing"
	"casino_common/common/service/countService"
)

func TestPushCountData(t *testing.T) {
	pushService.PoolInit("192.168.199.155:2801")
	//pushService.PushUserData(11750)
	row := countService.T_game_log{
		UserId: 123,
		IsWine:true,
		Bill: 9595,
	}
	t.Log(row.Insert())
}
