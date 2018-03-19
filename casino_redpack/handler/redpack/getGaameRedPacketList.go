package redpack

import (
	"casino_redpack/modules"
)

//获取房间列表
func GetGaameRedPacketListHandler(ctx *modules.Context) {
	//user := ctx.Session.Get("user")
	list := `{
	"code": 1,
	"message": "success",
	"request": [{
		"id": 1137,
		"type": 1,
		"money": "50.00",
		"all_membey": 5,
		"has_member": 4,
		"tail_number": 0,
		"nickname": "zhujimeizu",
		"headimgurl": "\/static\/userpic\/1190.jpg"
	}, {
		"id": 1138,
		"type": 1,
		"money": "100.00",
		"all_membey": 5,
		"has_member": 4,
		"tail_number": 0,
		"nickname": "\u4f60\u76f8\u8c8c\u5802\u5802",
		"headimgurl": "\/static\/userpic\/536.jpg"
	}]
}`

	ctx.Write([]byte(list))
}
