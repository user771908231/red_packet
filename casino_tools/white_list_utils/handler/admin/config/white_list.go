package config

import (
	"casino_tools/white_list_utils/modules"
	"casino_common/common/service/whiteListService"
	"fmt"
	"casino_common/common/consts"
	"casino_common/utils/redisUtils"
	"casino_common/common/log"
)

//配置
func WhiteListHandler(ctx *modules.Context) {
	game_id := ctx.QueryInt("gid")
	if game_id == 0 {
		game_id = 27
	}

	list := whiteListService.GetWhiteListByGid(int32(game_id))

	ctx.Data["gid"] = game_id
	ctx.Data["list"] = list

	ctx.Data["list_json"],_ = ctx.JSONString(list)

	ctx.HTML(200, "admin/config/white_list/list")
}

//
func WhiteListPostHandler(ctx *modules.Context) {
	game_id := ctx.QueryInt("gid")
	data := ctx.Query("data")

	rkey := fmt.Sprintf("%s_%02d", consts.RKEY_GAME_WHITE_LIST, game_id)

	err := redisUtils.Set(rkey, data)

	if err != nil {
		log.T("设置白名单失败 %s %s err:%v", rkey, data, err)
		ctx.Ajax(-1, "白名单设置失败！", nil)
	}else {
		log.T("设置白名单成功 %s %s", rkey, data)
		ctx.Ajax(1, "白名单设置成功！", nil)
	}
}

