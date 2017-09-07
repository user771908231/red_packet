package game

import (
	"casino_common/utils/testUtils"
	"casino_testtools/modules"
	"fmt"
	"encoding/json"
)

func GameTest(ctx *modules.Context) {
	desk_pwd := ctx.Query("desk_pwd")
	room_type := 16
	if ctx.QueryInt("room_type") == 15 {
		room_type = 15
	}

	ctx.Data["desk_pwd"] = desk_pwd
	ctx.Data["room_type"] = room_type
	ctx.HTML(200, "game/game")
}

func GameEdit(ctx *modules.Context) {
	gameId := ctx.QueryInt("gameid")
	deskPws := ctx.Query("deskpwd")

	if gameId <= 0 || deskPws == "" {
		ctx.Ajax(-1, "表单数据异常", nil)
		return
	}

	json_str := ctx.Query("data")
	fmt.Println(json_str)

	pokers := map[uint32][]int{}
	err := json.Unmarshal([]byte(json_str), &pokers)

	if err != nil {
		ctx.Ajax(-2, "表单数据异常", nil)
		return
	}

	fmt.Println(err, pokers)

	err = testUtils.SetDeskPreSendPokers(int32(gameId), deskPws, pokers)

	if err != nil {
		ctx.Ajax(-3, "调牌失败！Err:"+err.Error(), nil)
		return
	}

	ctx.Ajax(1, "调牌成功！", nil)
}
